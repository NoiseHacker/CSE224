package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	pb "numlog/numlog"
	"os"
	"time"

	"google.golang.org/protobuf/proto"
)

func logRequest(logentry *pb.NumberEntry) {

	numlogwire, err := proto.Marshal(logentry)
	if err != nil {
		log.Fatal(err)
	}

	logfile, err2 := os.OpenFile("log.db", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err2 != nil {
		log.Fatal(err2)
	}

	log.Printf("Writing out a log entry of size: %d\n", uint32(len(numlogwire)))
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(len(numlogwire)))

	logfile.Write(buf)
	logfile.Write(numlogwire)

	logfile.Close()
}

func main() {

	ln, err := net.Listen("tcp", "localhost:7777")
	if err != nil {
		log.Fatalf("failed to listen on port 7777: %v", err)
	}
	defer ln.Close()
	log.Println("Server listening on localhost:7777")

	// Accept a single connection.
	conn, err := ln.Accept()
	if err != nil {
		log.Fatalf("failed to accept connection: %v", err)
	}
	defer conn.Close()
	log.Printf("Client connected from %s\n", conn.RemoteAddr())

	// Read int32 values until the client closes the connection.
	for {
		var num int32
		// binary.Read will read exactly 4 bytes and interpret them little-endian
		if err := binary.Read(conn, binary.LittleEndian, &num); err != nil {
			if err == io.EOF {
				log.Println("Client closed the connection")
				break
			}
			log.Fatalf("failed to read from connection: %v", err)
		}

		numlog := &pb.NumberEntry{Ts: time.Now().Unix(), Number: num}
		logRequest(numlog)

		fmt.Printf("Received: %d\n", num)
	}

	log.Println("Server exiting")
}
