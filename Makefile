.PHONY: all clean run_put run_get

all:
	cd proto && make
	go build -o kv && go build -o kv-go-grpc ./plugin-go-grpc

clean:
	rm -f kv kv-go-grpc kv_hello

run_put:
	KV_PLUGIN="./kv-go-grpc" ./kv put hello world

run_get:
	KV_PLUGIN="./kv-go-grpc" ./kv get hello

