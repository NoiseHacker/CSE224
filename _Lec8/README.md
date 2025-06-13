

# to build:

First install protoc and the protoc-gen-go tools as per the gRPC docs

$ protoc --go_out=. --go_opt=paths=source_relative numlog/numlog.proto

# to run the protocol buffers server:

$ go run loggingnumserver/main.go

# to run the log printer:

$ go run logprinter/main.go
