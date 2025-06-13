package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"tritontube/internal/proto"
	"tritontube/internal/storage"

	"google.golang.org/grpc"
)

func main() {
	host := flag.String("host", "localhost", "Host address for the server")
	port := flag.Int("port", 8090, "Port number for the server")
	flag.Parse()

	// Validate arguments
	if *port <= 0 {
		panic("Error: Port number must be positive")
	}

	if flag.NArg() < 1 {
		fmt.Println("Usage: storage [OPTIONS] <baseDir>")
		fmt.Println("Error: Base directory argument is required")
		return
	}
	baseDir := flag.Arg(0)

	fmt.Println("Starting storage server...")
	fmt.Printf("Host: %s\n", *host)
	fmt.Printf("Port: %d\n", *port)
	fmt.Printf("Base Directory: %s\n", baseDir)

	addr := fmt.Sprintf("%s:%d", *host, *port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	grpcServer := grpc.NewServer()
	handler := storage.NewStorageHandler(baseDir)
	proto.RegisterStorageServiceServer(grpcServer, handler)

	log.Printf("Storage server listening at %s, baseDir: %s", addr, baseDir)
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
