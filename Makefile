all: protoc client server

protoc:
	@echo "Generating Go files"
	cd src/proto && protoc --go_out=plugins=grpc:. *.proto

server:
	@echo "Building server"
	go build -o server \
		github.com/pramonow/go-grpc-server-streaming-example/src/server

client:
	@echo "Building client"
	go build -o client \
		github.com/pramonow/go-grpc-server-streaming-example/src/client

.PHONY: client server protoc