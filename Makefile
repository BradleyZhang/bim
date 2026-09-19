.PHONY: build server client build-server build-client

build: build-server build-client

build-server:
	mkdir -p bin
	go -C server build -o ../bin/bim-server

build-client:
	mkdir -p bin
	go -C client build -o ../bin/bim-client

server: build-server
	./bin/bim-server

client: build-client
	./bin/bim-client