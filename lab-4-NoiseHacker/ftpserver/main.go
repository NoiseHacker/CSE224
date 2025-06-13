package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := flag.Int("port", 3333, "Port to accept connections on")
	host := flag.String("host", "127.0.0.1", "Host to bind to")
	delim := flag.String("delimiter", "\r\n", "Delimiter that separates commands")
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatalf("Usage: %s [--host HOST] [--port PORT] <directory>", os.Args[0])
	}
	dir := flag.Arg(0)

	log.Printf("%s:%d", *host, *port)
	address := *host + ":" + strconv.Itoa(*port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen to client: %v", err)
	}
	defer listener.Close()
	log.Printf("Server started on %s, serving files from %s", address, dir)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn, *delim, dir)
	}
}

func handleConnection(conn net.Conn, delim string, dir string) {
	log.Println("Accepted new connection.")
	defer conn.Close()
	defer log.Println("Closed connection.")

	// copy from lab3
	reader := bufio.NewReader(conn)      // <STOR|RETR> <filename> [rate] "\r\n"
	line, err := reader.ReadString('\n') // read line by line
	if err != nil {
		if err.Error() != "EOF" {
			log.Printf("Read error: %v", err)
		}
		return
	}
	line = strings.TrimRight(line, delim)
	parts := strings.Split(line, " ") // seperate by space
	if len(parts) < 2 {
		log.Println("Invalid command" + line)
    	return
	}

	cmd := parts[0]
	filename := filepath.Base(parts[1]) // extract the file name
	file_path := filepath.Join(dir, filename)

	if cmd == "STOR" {
		handleSTOR(reader, file_path)
	} else if cmd == "RETR" {
		rate := int64(0)
		if len(parts) > 2 { // with rate limit
			rate, _ = strconv.ParseInt(parts[2], 10, 64)
		}
		handleRETR(conn, file_path, rate)
	} else {
		log.Printf("Unknown command, Please use <STOR|RETR> in client side.")
	}
}

func handleSTOR(reader *bufio.Reader, file_path string) {
	// Server receives the file contents,
	// then stores the data as the specified filename on its local storage.
	file, err := os.Create(file_path)
	if err != nil {
		log.Printf("Create file error: %v", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		log.Println(err)
	}
}

func handleRETR(conn net.Conn, file_path string, rate int64) {
	// Return the contents of file the client issued,
	// Close TCP connection
	file, err := os.Open(file_path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	chunkSize := int64(4096)         // 4096 bytes
	chunk := make([]byte, chunkSize) // 4KB per chunk

	if rate == 0 { // no rate limit
		_, err = io.Copy(conn, file)
		if err != nil {
			log.Println(err)
		}
		return
	}
	for {
		n, err := file.Read(chunk)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %v", err)
			}
			break
		}
		_, err = conn.Write(chunk[:n])
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
		time.Sleep(time.Second * time.Duration(chunkSize*8) / time.Duration(rate)) // control the rate limits
	}
}
