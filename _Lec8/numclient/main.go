package main

import (
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:7777")
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to server at localhost:7777")

	// Number of random values to send
	r := rand.New(rand.NewSource(42))
	const count = 10

	for i := 0; i < count; i++ {
		var num int32 = int32(r.Int())

		if err := binary.Write(conn, binary.LittleEndian, num); err != nil {
			log.Fatalf("failed to write to connection: %v", err)
		}
		log.Printf("Sent: %d", num)
		time.Sleep(500 * time.Millisecond) // slight pause
	}

	log.Println("All numbers sent; closing connection")
}
