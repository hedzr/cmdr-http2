# cmdr-http2

A [`cmdr`](https://github.com/hedzr/cmdr) demo app.
`cmdr-http2` implements a http2 server with full daemon supports and graceful shutdown.

```bash

# clone and init
git clone https://github.com/hedzr/cmdr-http2.git
cd cmdr-http2
go mod download

# run server
go run cli/main.go server run &

# run client and make an request
go run cli/main.go h2

# or via curl
curl -k https://localhost:5151/

```