.PHONY: all dep

all:
	protoc --go_out=.. --go-grpc_out=.. kv.proto

dep:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

clean:
	rm ./*pb.go