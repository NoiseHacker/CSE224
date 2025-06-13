package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	port := flag.Int("port", 3333, "Port to accept connections on")
	host := flag.String("host", "127.0.0.1", "Host to bind to")
	flag.Parse()

	if len(flag.Args()) < 2 { // instead of os.Args
		log.Fatalf("Usage: %s <STOR|RETR> <filename> [rate]", os.Args[0])
	}

	log.Printf("Connecting to %s on port %d", *host, *port)
	address := *host + ":" + strconv.Itoa(*port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	cmd := flag.Arg(0)
	file_path := flag.Arg(1)
	rate := int64(0)
	if len(flag.Args()) > 2 { // include "rate" argument
		rate, err = strconv.ParseInt(flag.Arg(2), 10, 64)
		if err != nil {
			log.Panicln(err)
		}
	}

	if cmd == "STOR" {
		handleSTOR(conn, file_path)
	} else if cmd == "RETR" {
		handleRETR(conn, file_path, rate)
	} else {
		log.Printf("Unknown command, Please use <STOR|RETR> instead.")
	}
}

func handleSTOR(conn net.Conn, file_path string) {
	// The client issues a STOR command and provides the local file in disk.
	// The <CRLF> delimiter signifies the end of the STOR command.
	// Then, the client provides the contents of the file and closes the connection.
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatalf("Open file error: %v", err)
	}
	defer file.Close()

	cmd := "STOR " + filepath.Base(file_path) + "\r\n" // <CRLF>

	_, err = conn.Write([]byte(cmd))
	if err != nil {
		log.Fatalln(err)
	}

	written, err := io.Copy(conn, file)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Upload completed (%d bytes)", written)
}

func handleRETR(conn net.Conn, file_path string, rate int64) {
	// The client connects to the server, issues a RETR command,
	// then downloads the filename in server directory.
	cmd := "RETR " + filepath.Base(file_path)
	if rate > 0 {
		cmd += " " + strconv.Itoa(int(rate))
	}
	cmd += "\r\n"

	_, err := conn.Write([]byte(cmd))
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Create(filepath.Base(file_path))
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	written, err := io.Copy(file, conn)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Downloaded completed (%d bytes)", written)
}
