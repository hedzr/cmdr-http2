/*
 * Copyright © 2019 Hedzr Yeh.
 */

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/hedzr/cmdr/plugin/daemon/impl"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort = 5151
)

type (
	daemonImpl struct {
		appTag      string
		certManager *autocert.Manager
		mux         *http.ServeMux
	}
)

//
//
//

// newDaemon creates an `daemon.Daemon` object
func newDaemon() daemon.Daemon {
	return &daemonImpl{}
}

// func OnBuildCmd(root *cmdr.RootCommand) {
// 	cmdr.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {
//
// 		// app.server.port
// 		if cx := cmdr.FindSubCommand("server", &root.Command); cx != nil {
// 			// logrus.Debugf("`server` command found")
// 			opt := cmdr.NewCmdFrom(cx)
// 			if flg := cmdr.FindFlag("port", cx); flg != nil {
// 				flg.DefaultValue = defaultPort
//
// 			} else {
// 				opt.NewFlag(cmdr.OptFlagTypeInt).
// 					Titles("p", "port").
// 					Description("the port to listen.", "").
// 					Group("").
// 					DefaultValue(defaultPort, "PORT")
// 			}
// 		}
// 	})
// }

//
//
//

func (d *daemonImpl) OnInstall(ctx *impl.Context, cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("%s daemon OnInstall", cmd.GetRoot().AppName)
	return
}

func (d *daemonImpl) OnUninstall(ctx *impl.Context, cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("%s daemon OnUninstall", cmd.GetRoot().AppName)
	return
	// panic("implement me")
}

func (d *daemonImpl) OnStatus(ctx *impl.Context, cmd *cmdr.Command, p *os.Process) (err error) {
	fmt.Printf("%s v%v\n", cmd.GetRoot().AppName, cmd.GetRoot().Version)
	fmt.Printf("PID=%v\nLOG=%v\n", ctx.PidFileName, ctx.LogFileName)
	return
}

func (d *daemonImpl) OnReload() {
	logrus.Debugf("%s daemon OnReload", d.appTag)
}

func (d *daemonImpl) OnStop(cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("%s daemon OnStop", cmd.GetRoot().AppName)
	return
}

func (d *daemonImpl) OnHotReload(ctx *impl.Context) (err error) {
	logrus.Debugf("%s daemon OnHotReload, pid = %v, ppid = %v, ctx = %v", d.appTag, os.Getpid(), os.Getppid(), ctx)
	return
}

func (d *daemonImpl) domains(topDomains ...string) (domainList []string) {
	for _, top := range cmdr.GetStringSliceR("server.domains", topDomains...) {
		domainList = append(domainList, top)
		for _, s := range []string{"aurora", "api", "home", "res"} {
			domainList = append(domainList, fmt.Sprintf("%s.%s", s, top))
		}
	}
	return
}

func (d *daemonImpl) checkAndEnableAutoCert() (tlsConfig *tls.Config) {
	tlsConfig = &tls.Config{}

	if cmdr.GetBoolR("server.autocert.enabled") {
		logrus.Debugf("...autocert enabled")
		d.certManager = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(d.domains("example.com")...), // 测试时使用的域名：example.com
			Cache:      autocert.DirCache("ci/certs"),
		}
		go func() {
			if err := http.ListenAndServe(":80", d.certManager.HTTPHandler(nil)); err != nil {
				logrus.Fatal("autocert tool listening on :80 failed.", err)
			}
		}()
		tlsConfig.GetCertificate = d.certManager.GetCertificate
	}

	return
}

func (d *daemonImpl) enableGracefulShutdown(srv *http.Server, stopCh, doneCh chan struct{}) {

	go func() {
		for {
			select {
			case <-stopCh:
				logrus.Debugf("...shutdown going on.")
				ctx, cancelFunc := context.WithTimeout(context.TODO(), 8*time.Second)
				defer cancelFunc()
				if err := srv.Shutdown(ctx); err != nil {
					logrus.Error("Shutdown failed: ", err)
				} else {
					logrus.Debugf("Shutdown ok.")
				}
				<-doneCh
				return
			}
		}
	}()

}

func (d *daemonImpl) buildRoutes(mux *http.ServeMux) (err error) {
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/", echoHandler)
	return
}

func (d *daemonImpl) OnRun(cmd *cmdr.Command, args []string, stopCh, doneCh chan struct{}, listener net.Listener) (err error) {
	d.appTag = cmd.GetRoot().AppName
	logrus.Debugf("%s daemon OnRun, pid = %v, ppid = %v", d.appTag, os.Getpid(), os.Getppid())

	port := cmdr.GetInt("app.server.port")
	if port == 0 {
		logrus.Fatal("port not defined")
	}

	// Tweak configuration values here.
	var (
		addr = fmt.Sprintf(":%d", port) // ":3300"
	)

	d.mux = http.NewServeMux()
	err = d.buildRoutes(d.mux)
	if err != nil {
		return
	}

	// Create a server on port 8000
	// Exactly how you would run an HTTP/1.1 server
	srv := &http.Server{
		Addr:              addr,
		Handler:           d.mux, // http.HandlerFunc(d.handle),
		TLSConfig:         d.checkAndEnableAutoCert(),
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	d.enableGracefulShutdown(srv, stopCh, doneCh)

	// TODO server push, ...
	// https://posener.github.io/http2/

	go func() {
		// Start the server with TLS, since we are running HTTP/2 it must be
		// run with TLS.
		// Exactly how you would run an HTTP/1.1 server with TLS connection.
		if srv.TLSConfig.GetCertificate == nil {
			logrus.Printf("Serving on https://0.0.0.0:%d ...", port)
			if cmdr.FileExists("ci/certs/server.cert") && cmdr.FileExists("ci/certs/server.key") {
				if err = d.serve(srv, listener, "ci/certs/server.cert", "ci/certs/server.key"); err != http.ErrServerClosed {
					logrus.Fatal(err)
				}
				logrus.Println("end")
			} else {
				logrus.Fatalf(`ci/certs/server.{cert,key} NOT FOUND under '%s'. You might generate its at command line:

	[ -d ci/certs ] || mkdir -p ci/certs
	openssl genrsa -out ci/certs/server.key 2048
	openssl req -new -x509 -key ci/certs/server.key -out ci/certs/server.cert -days 3650 -subj /CN=localhost

				`, cmdr.GetCurrentDir())
			}
		} else {
			logrus.Printf("Serving on https://0.0.0.0:%d with autocert...", port)
			if err = d.serve(srv, listener, "", ""); err != http.ErrServerClosed {
				logrus.Fatal(err)
			}
			logrus.Println("end")
		}
	}()

	// go worker(stopCh, doneCh)
	return
}

func (d *daemonImpl) serve(srv *http.Server, listener net.Listener, certFile, keyFile string) (err error) {
	// if srv.shuttingDown() {
	// 	return http.ErrServerClosed
	// }

	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	if listener == nil {
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			return err
		}
	}

	defer func() {
		h2listener.Close()
		logrus.Printf("h2listener closed, pid=%v", os.Getpid())
	}()

	h2listener = listener

	return srv.ServeTLS(h2listener, certFile, keyFile)
}

func (d *daemonImpl) worker(stopCh, doneCh chan struct{}) {
LOOP:
	for {
		time.Sleep(3 * time.Second) // this is work to be done by worker.
		select {
		case <-stopCh:
			break LOOP
		default:
			logrus.Debugf("%s running at %d", d.appTag, os.Getpid())
		}
	}
	doneCh <- struct{}{}
}

func (d *daemonImpl) handle(w http.ResponseWriter, r *http.Request) {
	// Log the request protocol
	log.Printf("Got connection: %s", r.Proto)
	// Send a message back to the client
	_, _ = w.Write([]byte("Hello"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Hello, world!\n")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, r.URL.Path)
}

const (
	// for http client
	activeTimeout       = 10 * time.Minute
	maxIdleConns        = 1000
	maxIdleConnsPerHost = 100

	// for http server
	idleTimeout       = 5 * time.Minute
	readHeaderTimeout = 1 * time.Second
	writeTimeout      = 10 * time.Second
	maxHeaderBytes    = http.DefaultMaxHeaderBytes
	shutdownTimeout   = 30 * time.Second
)
