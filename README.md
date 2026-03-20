
install atlas (local)
install make
### Generate Code

goctl rpc protoc greeter.proto --go_out=. --go-grpc_out=. --zrpc_out=.


brew install yq

// create file pb for gateway
protoc --descriptor_set_out=dropshipbe.pb dropshipbe.proto# dropship-be
