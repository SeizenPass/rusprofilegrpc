genG:
	protoc -I . \
        --go_out . \
        --go_opt paths=source_relative \
        --go-grpc_out . \
        --go-grpc_opt paths=source_relative \
        proto/rusprofile.proto
genW:
	protoc -I . \
	  --grpc-gateway_out . \
	  --grpc-gateway_opt logtostderr=true \
	  --grpc-gateway_opt paths=source_relative \
	  --grpc-gateway_opt grpc_api_configuration=proto/search.yaml \
	  proto/rusprofile.proto
build: genG genW
clean:
	rm proto/*.go
runs:
	go run server/server.go
runc:
	go run client/client.go