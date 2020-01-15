// Copyright © 2020 Hedzr Yeh.

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/hedzr/cmdr"
	tls2 "github.com/hedzr/cmdr-http2/cli/server/tls"
	"github.com/hedzr/cmdr/plugin/daemon"
	"github.com/hedzr/cmdr/plugin/daemon/impl"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/host"
	"github.com/kataras/iris/v12/core/netutil"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

const (
	defaultPort = 1379

	typeDefault muxType = iota
	typeGin
	typeIris
)

type (
	muxType int

	daemonImpl struct {
		appTag      string
		certManager *autocert.Manager
		Type        muxType
		mux         *http.ServeMux
		router      *gin.Engine
		irisApp     *iris.Application
	}
)

//
//
//

// newDaemon creates an `daemon.Daemon` object
func newDaemon() daemon.Daemon {
	return &daemonImpl{Type: typeIris}
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

func (d *daemonImpl) domains() (domainList []string) {
	for _, top := range cmdr.GetStringSliceR("server.autocert.domains", "example.com") {
		domainList = append(domainList, top)
		for _, s := range cmdr.GetStringSliceR("server.autocert.second-level-domains", "aurora", "api", "home", "res") {
			domainList = append(domainList, fmt.Sprintf("%s.%s", s, top))
		}
	}
	return
}

func (d *daemonImpl) checkAndEnableAutoCert(config *tls2.CmdrTLSConfig) (tlsConfig *tls.Config) {
	tlsConfig = &tls.Config{}

	if config.IsServerCertValid() {
		tlsConfig = config.ToServerTLSConfig()
	}

	if cmdr.GetBoolR("server.autocert.enabled") {
		logrus.Debugf("...autocert enabled")
		d.certManager = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(d.domains()...), // 测试时使用的域名：example.com
			Cache:      autocert.DirCache(cmdr.GetStringR("server.autocert.dir-cache", "ci/certs")),
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

func (d *daemonImpl) getHandler() http.Handler {
	switch d.Type {
	case typeGin:
		return d.router
	case typeIris:
		return d.irisApp
	default:
		return d.mux
	}
}

func (d *daemonImpl) OnRun(cmd *cmdr.Command, args []string, stopCh, doneCh chan struct{}, listener net.Listener) (err error) {
	d.appTag = cmd.GetRoot().AppName
	logrus.Debugf("%s daemon OnRun, pid = %v, ppid = %v", d.appTag, os.Getpid(), os.Getppid())

	port := cmdr.GetIntR("cmdr-http2.server.port")

	// Tweak configuration values here.
	var (
		config    = tls2.NewCmdrTLSConfig("cmdr-http2.server.tls", "server.start")
		tlsConfig = d.checkAndEnableAutoCert(config)
	)

	logrus.Tracef("used config file: %v", cmdr.GetUsedConfigFile())
	logrus.Tracef("logger level: %v", logrus.GetLevel())

	if config.IsServerCertValid() || tlsConfig.GetCertificate == nil {
		port = cmdr.GetIntR("cmdr-http2.server.ports.tls")
	}

	switch d.Type {
	case typeGin:
		// d.router = gin.Default()
		gin.ForceConsoleColor()
		d.router = gin.New()
		d.router.Use(gin.Logger())
		d.router.Use(gin.Recovery())
		// d.router.GET("/benchmark", MyBenchLogger(), benchEndpoint)
		err = d.buildGinRoutes(d.router)

	case typeIris:
		d.irisApp = iris.New()
		d.irisApp.Logger().SetLevel("debug")
		d.irisApp.Use(recover.New())
		d.irisApp.Use(logger.New())
		err = d.buildIrisRoutes(d.irisApp)

	default:
		d.mux = http.NewServeMux()
		err = d.buildRoutes(d.mux)
	}

	if err != nil {
		return
	}

	if port == 0 {
		logrus.Fatal("port not defined")
	}
	addr := fmt.Sprintf(":%d", port) // ":3300"

	// Create a server on port 8000
	// Exactly how you would run an HTTP/1.1 server
	srv := &http.Server{
		Addr:              addr,
		Handler:           d.getHandler(), // d.mux, // http.HandlerFunc(d.handle),
		TLSConfig:         tlsConfig,
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
		if config.IsServerCertValid() || srv.TLSConfig.GetCertificate == nil {
			logrus.Printf("Serving on https://0.0.0.0:%d with HTTPS...", port)
			// if cmdr.FileExists("ci/certs/server.cert") && cmdr.FileExists("ci/certs/server.key") {
			if err = d.serve(srv, listener, config.Cert, config.Key); err != http.ErrServerClosed {
				logrus.Fatalf("listen: %s\n", err)
			}
			// if err = d.serve(srv, listener, "ci/certs/server.cert", "ci/certs/server.key"); err != http.ErrServerClosed {
			// 	logrus.Fatal(err)
			// }
			logrus.Println("end")
			// 		} else {
			// 			logrus.Fatalf(`ci/certs/server.{cert,key} NOT FOUND under '%s'. You might generate its at command line:
			//
			// [ -d ci/certs ] || mkdir -p ci/certs
			// openssl genrsa -out ci/certs/server.key 2048
			// openssl req -new -x509 -key ci/certs/server.key -out ci/certs/server.cert -days 3650 -subj /CN=localhost
			//
			// 			`, cmdr.GetCurrentDir())
			// 		}
		} else {
			logrus.Printf("Serving on https://0.0.0.0:%d with HTTP...", port)
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

	switch d.Type {
	case typeIris:
		return d.irisApp.Run(iris.Raw(func() error {
			su := d.irisApp.NewHost(srv)
			if netutil.IsTLS(su.Server) {
				h2listener = tls.NewListener(h2listener, su.Server.TLSConfig)
				su.Configure(func(su *host.Supervisor) {
					rs := reflect.ValueOf(su).Elem()
					// rf := rs.FieldByName("manuallyTLS")
					rf := rs.Field(2)
					// rf can't be read or set.
					rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
					// Now rf can be read and set.

					// su.manuallyTLS = true
					i := true
					ri := reflect.ValueOf(&i).Elem() // i, but writeable
					rf.Set(ri)
				})
			}
			err = su.Serve(h2listener)
			return err
		}), iris.WithoutServerError(iris.ErrServerClosed))
	default:
		return srv.ServeTLS(h2listener, certFile, keyFile)
	}
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

// https://revel.github.io/
// https://github.com/revel/revel

func (d *daemonImpl) buildIrisRoutes(app *iris.Application) (err error) {
	// https://iris-go.com/start/
	// https://github.com/kataras/iris
	//
	// https://www.slant.co/topics/1412/~best-web-frameworks-for-go

	app.Get("/", func(c iris.Context) {
		_, _ = c.JSON(iris.Map{"message": "Hello Iris!"})
	})
	app.Get("/ping", func(ctx iris.Context) {
		_, _ = ctx.WriteString("pong")
	})
	// Resource: http://localhost:1380
	app.Handle("GET", "/welcome", func(ctx iris.Context) {
		_, _ = ctx.HTML("<h1>Welcome</h1>")
	})

	// app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))

	app.Get("/s/:path", d.echoIrisHandler)
	return
}

func (d *daemonImpl) echoIrisHandler(ctx iris.Context) {
	p := ctx.Params().GetString("path")
	if p == "zero" {
		d0(8, 0) // raise a 0-divide panic and it will be recovered by http.Conn.serve(ctx)
	}
	_, _ = ctx.WriteString(p)
}

func (d *daemonImpl) buildGinRoutes(mux *gin.Engine) (err error) {
	// https://github.com/gin-gonic/gin
	// https://github.com/gin-contrib/multitemplate
	// https://github.com/gin-contrib
	//
	// https://www.mindinventory.com/blog/top-web-frameworks-for-development-golang/

	mux.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	mux.GET("/hello", helloGinHandler)

	mux.GET("/s/*action", echoGinHandler)

	// mux.Static("/assets", "./assets")
	mux.StaticFS("/public", http.Dir("./public"))
	mux.StaticFile("/favicon.ico", "./resources/favicon.ico")
	
	return
}

func helloGinHandler(c *gin.Context) {
	_, _ = io.WriteString(c.Writer, "Hello, world!\n")
}

func echoGinHandler(c *gin.Context) {
	action := c.Param("action")
	if action == "/zero" {
		d0(8, 0) // raise a 0-divide panic and it will be recovered by http.Conn.serve(ctx)
	}
	_, _ = io.WriteString(c.Writer, action)
}

func (d *daemonImpl) buildRoutes(mux *http.ServeMux) (err error) {
	// https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
	//

	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/", echoHandler)
	return
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
	if r.URL.Path == "/zero" {
		d0(8, 0) // raise a 0-divide panic and it will be recovered by http.Conn.serve(ctx)
	}
}

func d0(a, b int) int {
	return a / b
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
