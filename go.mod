module github.com/hedzr/cmdr-http2

go 1.12

// replace github.com/hedzr/cmdr => ../cmdr

// replace github.com/hedzr/logex v0.0.0 => ../logex

// replace github.com/hedzr/pools v0.0.0 => ../pools

// replace github.com/hedzr/errors v0.0.0 => ../errors

require (
	github.com/hedzr/cmdr v1.6.12
	github.com/hedzr/errors v1.1.11
	github.com/hedzr/logex v1.1.3
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3
	gopkg.in/yaml.v2 v2.2.2
)
