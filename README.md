# cmdr-http2

A [`cmdr`](https://github.com/hedzr/cmdr) demo app.
`cmdr-http2` implements a http2 server with full daemon supports and graceful shutdown.

NOTE: need cmdr v1.6.25+ [20200204]

```bash

# clone and init
git clone https://github.com/hedzr/cmdr-http2.git
cd cmdr-http2
go mod download

# run server
go run ./cli/ server run &

# run client and make an request
go run ./cli/ h2

# or via curl
curl -k https://localhost:5151/

#
# Build the binary
#
go build -o bin/cli ./cli/
# or:
make build

# Shell prompt mode
$ go run ./cli/ shell
>>> --help
>>> quit
# type <space> to get auto-completion tip
# type sub-commands
```


### Plugins for `cmdr` System

#### `sample`

sample plugin give a example to howto modify cmdr daemon plugin `server` `start` command at an appropriate time.



#### `trace`

trace plugin adds a `trace` option to cmdr system.



#### `shell`

enable shell prompt mode inside app.



### `cmdr-http2` Shell Prompt Mode:

![image](https://user-images.githubusercontent.com/12786150/71587009-11436500-2b57-11ea-890d-a60989a09248.png)

#### Shell prompt mode

the feature is powered by [c-bata/go-prompt](https://github.com/c-bata/go-prompt).




## New H2 Server since v1.3.5

- supports autocert, std TLS, ...
- supports graceful-shutdown
- supports hot-reload
- demostrates howto use various of Go Web Frameworks: std go/http, iris, gin, ...


### create certificates for testing

This command will generate CA, server, client certificates and wrote them into `./ci/certs` for h2 server loading.

```bash
bin/cmdr-http2 server certs create
bin/cmdr-http2 server certs create --help
bin/cmdr-http2 server certs --help
```

See the source codes for more information.

### configuration file(s)

In the project root directory, `cmdr-http2.yml` will be loaded as main config file, this depends on `cmdr` config file searching algorithm. See also `cmdr.WithPredefinedLocations(locations...)`.





## LICENSE

MIT
