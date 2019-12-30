# cmdr-http2

A [`cmdr`](https://github.com/hedzr/cmdr) demo app.
`cmdr-http2` implements a http2 server with full daemon supports and graceful shutdown.

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





### cmdr-http2 shell prompt mode:

![image](https://user-images.githubusercontent.com/12786150/71587009-11436500-2b57-11ea-890d-a60989a09248.png)

#### Shell prompt mode

the feature is powered by [c-bata/go-prompt](https://github.com/c-bata/go-prompt).





# LICENSE

MIT
