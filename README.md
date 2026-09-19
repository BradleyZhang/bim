A simple instant messaging system built with Go.

## Run
### Using Makefile
Start server: `make server`

Start client `make client`

### Using go directly
Start server `go run ./server`

Start client `go run ./client`

## Build
### Using Makefile

```bash
make build #build server and client

make build-server
make build-client
```

### Using Go directly

Server
```bash
go -C server build -o ../bin/bim-server
```

Client
```bash
go -C client build -o ../bin/bim-client
```