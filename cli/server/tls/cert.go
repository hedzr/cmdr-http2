// Copyright © 2020 Hedzr Yeh.

package tls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path"
	"time"
)

func CaCreate(cmd *cmdr.Command, args []string) (err error) {
	return
}

var pkixName pkix.Name = pkix.Name{
	Country:            []string{"CN"},
	Locality:           []string{"CQ"},
	Province:           []string{"Chongqing"},
	StreetAddress:      []string{},
	PostalCode:         []string{},
	Organization:       []string{"Hedzr Studio."},
	OrganizationalUnit: []string{},
	CommonName:         "Root CA",
}

var outputDirs = []string{"./ci/certs"}

const (
	rootKeyFileName        = "root.key"
	rootCertFileName       = "root.pem"
	rootCertDbgFileName    = "root.debug.crt"
	leafKeyFileName        = "cert.key"
	leafCertFileName       = "cert.pem"
	leafCertDbgFileName    = "cert.debug.crt"
	leafCertBundleFileName = "server-bundle.pem"
	clientKeyFileName      = "client.key"
	clientCertFileName     = "client.pem"
	clientCertDbgFileName  = "client.debug.crt"

	rootCACommonName = "Root CA"
	leafCommonName   = "localhost"
	clientCommonName = "demo1"
)

func CertCreate(cmd *cmdr.Command, args []string) (err error) {
	var prefix = "server.certs"

	var notBefore time.Time
	validFrom := cmdr.GetStringRP(prefix, "create.start-date")
	if len(validFrom) == 0 {
		notBefore = time.Now().UTC()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", validFrom)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse creation date: %s\n", err)
			os.Exit(1)
		}
	}

	validFor := cmdr.GetDurationRP(prefix, "create.valid-for")
	notAfter := notBefore.Add(validFor)

	hosts := cmdr.GetStringSliceRP(prefix, "create.host")

	outputDir := ""
	for _, dir := range outputDirs {
		err = cmdr.EnsureDir(dir)
		if err != nil {
			panic(err)
		}
		outputDir = dir
	}

	var (
		// derBytes     []byte
		// serialNumber *big.Int
		rootKey      *ecdsa.PrivateKey
		rootTemplate *x509.Certificate
		caCertBytes  []byte
	)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)

	rootTemplate, rootKey, caCertBytes, err = newCaCerts(outputDir, notBefore, notAfter, serialNumberLimit)
	if err != nil {
		panic(err)
	}
	err = newLeafCerts(outputDir, notBefore, notAfter, serialNumberLimit, rootTemplate, rootKey, caCertBytes, hosts)
	if err != nil {
		panic(err)
	}
	err = newClientCerts(outputDir, notBefore, notAfter, rootTemplate, rootKey)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintf(os.Stdout, `Successfully generated certificates! Here's what you generated.
# Root CA
%v
	The private key for the root Certificate Authority. Keep this private.
%v
	The public key for the root Certificate Authority. Clients should load the
	certificate in this file to connect to the server.
%v
	Debug information about the generated certificate.

# Leaf Certificate - Use these to serve TLS traffic.
%v
	Private key (PEM-encoded) for terminating TLS traffic on the server.
%v
	Public key for terminating TLS traffic on the server.
%v
	Debug information about the generated certificate

# Client Certificate - You probably don't need these.
%v: Secret key for TLS client authentication
%v: Public key for TLS client authentication

`,
		rootKeyFileName, rootCertFileName, rootCertDbgFileName,
		leafKeyFileName, leafCertFileName, leafCertDbgFileName,
		clientKeyFileName, clientCertFileName,
	)
	return
}

func newCaCerts(outputDir string, notBefore, notAfter time.Time, serialNumberLimit *big.Int) (rootTemplate *x509.Certificate, rootKey *ecdsa.PrivateKey, caCertBytes []byte, err error) {

	var (
		serialNumber *big.Int
		derBytes     []byte
		caKeyPath    string
		caPath       string
	)

	caKeyPath = path.Join(outputDir, rootKeyFileName)
	caPath = path.Join(outputDir, rootCertFileName)

	if cmdr.FileExists(caKeyPath) && cmdr.FileExists(caPath) {
		logrus.Infof("ignore recreating certs: %v, %v", caKeyPath, caPath)
		return // exists, ignore creating
	}

	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logrus.Fatalf("failed to generate serial number: %s", err)
	}

	rootKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	keyToFile(caKeyPath, rootKey)

	pkixName.CommonName = rootCACommonName
	rootTemplate = &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkixName,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err = x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		panic(err)
	}
	debugCertToFile(path.Join(outputDir, rootCertDbgFileName), derBytes)
	certToFile(caPath, derBytes)

	caCertBytes = derBytes

	return
}

func newLeafCerts(outputDir string, notBefore, notAfter time.Time, serialNumberLimit *big.Int, rootTemplate *x509.Certificate, rootKey *ecdsa.PrivateKey, caCertBytes []byte, hosts []string) (err error) {

	// http.ListenAndServeTLS(":7252", "leaf.pem", "leaf.key", nil)

	var (
		serialNumber          *big.Int
		leafKey               *ecdsa.PrivateKey
		derBytes              []byte
		cKeyPath              string
		cPath, cbPath, caPath string
	)

	cKeyPath = path.Join(outputDir, leafKeyFileName)
	cPath = path.Join(outputDir, leafCertFileName)
	cbPath = path.Join(outputDir, leafCertBundleFileName)
	caPath = path.Join(outputDir, rootCertFileName)

	if cmdr.FileExists(cKeyPath) && cmdr.FileExists(cPath) {
		logrus.Infof("ignore recreating certs: %v, %v", cKeyPath, cPath)
		return // exists, ignore creating
	}

	leafKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	keyToFile(cKeyPath, leafKey)

	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logrus.Fatalf("failed to generate serial number: %s", err)
	}

	pkixName.CommonName = leafCommonName
	leafTemplate := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkixName,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			leafTemplate.IPAddresses = append(leafTemplate.IPAddresses, ip)
		} else {
			leafTemplate.DNSNames = append(leafTemplate.DNSNames, h)
		}
	}

	derBytes, err = x509.CreateCertificate(rand.Reader, &leafTemplate, rootTemplate, &leafKey.PublicKey, rootKey)
	if err != nil {
		panic(err)
	}
	debugCertToFile(path.Join(outputDir, leafCertDbgFileName), derBytes)
	certToFile(cPath, derBytes)

	certToFile(cbPath, derBytes)
	appendCertToFile(cbPath, caPath)

	return
}

func newClientCerts(outputDir string, notBefore, notAfter time.Time, rootTemplate *x509.Certificate, rootKey *ecdsa.PrivateKey) (err error) {

	var (
		serialNumber *big.Int
		clientKey    *ecdsa.PrivateKey
		derBytes     []byte
		cKeyPath     string
		cPath        string
	)

	cKeyPath = path.Join(outputDir, clientKeyFileName)
	cPath = path.Join(outputDir, clientCertFileName)

	if cmdr.FileExists(cKeyPath) && cmdr.FileExists(cPath) {
		logrus.Infof("ignore recreating certs: %v, %v", cKeyPath, cPath)
		return // exists, ignore creating
	}

	clientKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	keyToFile(cKeyPath, clientKey)

	serialNumber = new(big.Int).SetInt64(4)
	pkixName.CommonName = clientCommonName
	clientTemplate := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkixName,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	derBytes, err = x509.CreateCertificate(rand.Reader, &clientTemplate, rootTemplate, &clientKey.PublicKey, rootKey)
	if err != nil {
		panic(err)
	}
	debugCertToFile(path.Join(outputDir, clientCertDbgFileName), derBytes)
	certToFile(cPath, derBytes)

	return
}

// keyToFile writes a PEM serialization of |key| to a new file called
// |filename|.
func keyToFile(filename string, key *ecdsa.PrivateKey) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
		os.Exit(2)
	}
	if err := pem.Encode(file, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}); err != nil {
		panic(err)
	}
}

func certToFile(filename string, derBytes []byte) {
	certOut, err := os.Create(filename)
	if err != nil {
		logrus.Fatalf("failed to open %v for writing: %s", filename, err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		logrus.Fatalf("failed to write data to %v: %s", filename, err)
	}
	if err := certOut.Close(); err != nil {
		logrus.Fatalf("error closing %v: %s", filename, err)
	}
}

func appendCertToFile(filename, appendFilename string) {
	certOut, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		logrus.Fatalf("failed to open %v for writing: %s", filename, err)
	}

	var c []byte
	c, err = ioutil.ReadFile(appendFilename)
	if err != nil {
		logrus.Fatalf("failed to open %v for reading: %s", appendFilename, err)
	}

	_, err = certOut.Write(c)
	if err != nil {
		logrus.Fatalf("failed to write data to %v: %s", filename, err)
	}
	if err := certOut.Close(); err != nil {
		logrus.Fatalf("error closing %v: %s", filename, err)
	}
}

// debugCertToFile writes a PEM serialization and OpenSSL debugging dump of
// |derBytes| to a new file called |filename|.
func debugCertToFile(filename string, derBytes []byte) {
	cmd := exec.Command("openssl", "x509", "-text", "-inform", "DER")

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	cmd.Stdout = file
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	if _, err := stdin.Write(derBytes); err != nil {
		panic(err)
	}
	stdin.Close()
	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}
