A simple instant messaging system built with Go.

## Run

Server
```bash
go -C server build -o ../bin/bim
./bin/bim
```

Client
Connect to Server using Netcat
```bash
nc 127.0.0.1 8888
```