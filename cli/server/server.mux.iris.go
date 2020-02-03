package server

import (
	"crypto/tls"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/host"
	"github.com/kataras/iris/v12/core/netutil"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"net"
	"net/http"
	"reflect"
	"unsafe"
)

// https://revel.github.io/
// https://github.com/revel/revel

// func (d *daemonImpl) buildIrisRoutes(app *iris.Application) (err error) {
// 	// https://iris-go.com/start/
// 	// https://github.com/kataras/iris
// 	//
// 	// https://www.slant.co/topics/1412/~best-web-frameworks-for-go
//
// 	app.Get("/", func(c iris.Context) {
// 		_, _ = c.JSON(iris.Map{"message": "Hello Iris!"})
// 	})
// 	app.Get("/ping", func(ctx iris.Context) {
// 		_, _ = ctx.WriteString("pong")
// 	})
// 	// Resource: http://localhost:1380
// 	app.Handle("GET", "/welcome", func(ctx iris.Context) {
// 		_, _ = ctx.HTML("<h1>Welcome</h1>")
// 	})
//
// 	// app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
//
// 	app.Get("/s/:path", d.echoIrisHandler)
// 	return
// }
//
// func (d *daemonImpl) echoIrisHandler(ctx iris.Context) {
// 	p := ctx.Params().GetString("path")
// 	if p == "zero" {
// 		d0(8, 0) // raise a 0-divide panic and it will be recovered by http.Conn.serve(ctx)
// 	}
// 	_, _ = ctx.WriteString(p)
// }

func newIris() *irisImpl {
	d := &irisImpl{}
	d.init()
	return d
}

type irisImpl struct {
	irisApp *iris.Application
}

func (d *irisImpl) Handler() http.Handler {
	return d.irisApp
}

func (d *irisImpl) Serve(srv *http.Server, listener net.Listener, certFile, keyFile string) (err error) {
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
}

func (d *irisImpl) BuildRoutes() {
	// https://iris-go.com/start/
	// https://github.com/kataras/iris
	//
	// https://www.slant.co/topics/1412/~best-web-frameworks-for-go

	d.irisApp.Get("/", func(c iris.Context) {
		_, _ = c.JSON(iris.Map{"message": "Hello Iris!"})
	})
	d.irisApp.Get("/ping", func(ctx iris.Context) {
		_, _ = ctx.WriteString("pong")
	})
	// Resource: http://localhost:1380
	d.irisApp.Handle("GET", "/welcome", func(ctx iris.Context) {
		_, _ = ctx.HTML("<h1>Welcome</h1>")
	})

	// app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))

	d.irisApp.Get("/s/:path", d.echoIrisHandler)
}

func (d *irisImpl) init() {
	d.irisApp = iris.New()
	d.irisApp.Logger().SetLevel("debug")
	d.irisApp.Use(recover.New())
	d.irisApp.Use(logger.New())
}

func (d *irisImpl) echoIrisHandler(ctx iris.Context) {
	p := ctx.Params().GetString("path")
	if p == "zero" {
		d0(8, 0) // raise a 0-divide panic and it will be recovered by http.Conn.serve(ctx)
	}
	_, _ = ctx.WriteString(p)
}
