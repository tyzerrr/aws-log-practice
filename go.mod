module github.com/tyzerrr/aws-log-practice

go 1.26.2

tool (
	connectrpc.com/connect/cmd/protoc-gen-connect-go
	google.golang.org/protobuf/cmd/protoc-gen-go
)

require (
	connectrpc.com/connect v1.20.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
