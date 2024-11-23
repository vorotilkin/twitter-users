generate:
	protoc -I=./ --go_out=./ --go-grpc_out=./ ./users.proto