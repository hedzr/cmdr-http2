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
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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

// NewDaemon creates an `daemon.Daemon` object
func NewDaemon() daemon.Daemon {
	return &daemonImpl{}
}

func OnBuildCmd(root *cmdr.RootCommand) {
	cmdr.AddOnBeforeXrefBuilding(func(root *cmdr.RootCommand, args []string) {

		// app.server.port
		if cmd := cmdr.FindSubCommand("server", &root.Command); cmd != nil {

			// logrus.Debugf("`server` command found")

			opt := cmdr.NewCmdFrom(cmd)

			opt.NewFlag(cmdr.OptFlagTypeInt).
				Titles("p", "port").
				Description("the port to listen.", "").
				Group("").
				DefaultValue(5151, "PORT")

		}
	})
}

//
//
//

func (d *daemonImpl) OnInstall(cxt *daemon.Context, cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("%s daemon OnInstall", cmd.GetRoot().AppName)
	return
	// panic("implement me")
}

func (d *daemonImpl) OnUninstall(cxt *daemon.Context, cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("%s daemon OnUninstall", cmd.GetRoot().AppName)
	return
	// panic("implement me")
}

func (d *daemonImpl) OnStatus(cxt *daemon.Context, cmd *cmdr.Command, p *os.Process) (err error) {
	fmt.Printf("%s v%v\n", cmd.GetRoot().AppName, cmd.GetRoot().Version)
	fmt.Printf("PID=%v\nLOG=%v\n", cxt.PidFileName, cxt.LogFileName)
	return
}

func (d *daemonImpl) OnReload() {
	logrus.Debugf("%s daemon OnReload", d.appTag)
}

func (d *daemonImpl) OnStop(cmd *cmdr.Command, args []string) (err error) {
	logrus.Debugf("%s daemon OnStop", cmd.GetRoot().AppName)
	return
}

func (d *daemonImpl) domains(top string) (domainList []string) {
	domainList = append(domainList, top)
	for _, s := range []string{"aurora", "api", "home", "res"} {
		domainList = append(domainList, fmt.Sprintf("%s.%s", s, top))
	}
	return
}

func (d *daemonImpl) checkAndEnableAutoCert() (tlsConfig *tls.Config) {
	tlsConfig = &tls.Config{}

	if cmdr.GetBool("app.server.autocert.enabled") {
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

func (d *daemonImpl) OnRun(cmd *cmdr.Command, args []string, stopCh, doneCh chan struct{}) (err error) {
	d.appTag = cmd.GetRoot().AppName
	logrus.Debugf("%s daemon OnRun, pid = %v, ppid = %v", d.appTag, os.Getpid(), os.Getppid())

	port := cmdr.GetInt("app.server.port")
	if port == 0 {
		logrus.Fatal("port not defined")
	}

	d.mux = http.NewServeMux()
	err = d.buildRoutes(d.mux)
	if err != nil {
		return
	}

	// Create a server on port 8000
	// Exactly how you would run an HTTP/1.1 server
	srv := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   d.mux, // http.HandlerFunc(d.handle),
		TLSConfig: d.checkAndEnableAutoCert(),
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
			logrus.Fatal(srv.ListenAndServeTLS("ci/certs/server.cert", "ci/certs/server.key"))
		} else {
			logrus.Printf("Serving on https://0.0.0.0:%d with autocert...", port)
			logrus.Fatal(srv.ListenAndServeTLS("", ""))
		}
	}()

	// go worker(stopCh, doneCh)
	return
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

const idleTimeout = 5 * time.Minute
const activeTimeout = 10 * time.Minute
const maxIdleConns = 1000
const maxIdleConnsPerHost = 100
