package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)
// Client
func main1() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s host:port", os.Args[0])
	}
	service := os.Args[1]

	conn, err := net.Dial("tcp", service)
	checkError(err)

	result, err := readFully(conn)
	checkError(err)

	fmt.Println(string(result))
}

func readFully(conn net.Conn) ([]byte, error) {
	defer func() {
		_ = conn.Close()
	}()

	result := bytes.NewBuffer(nil)
	var buf [512]byte

	for {
		n, err := conn.Read(buf[:])
		result.Write(buf[:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return result.Bytes(), nil

}

// func checkError(err error) {
// 	if err != nil {
// 		log.Fatalf("Fatal error: %s", err.Error())
// 	}
// }
