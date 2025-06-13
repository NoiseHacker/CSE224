package main

import (
	"flag"
	"log"
	"net"
	"strconv"
	"strings"
	"bufio"
	"fmt"
)

func main() {
	// Command-line flags for host and port
	// name  value 	usage
	port := flag.Int("port", 3333, "Port to accept connections on")
	host := flag.String("host", "127.0.0.1", "Host to bind to")
	delim := flag.String("delimiter", "\r\n", "Delimiter that separates commands")
	flag.Parse()

	address := *host + ":" + strconv.Itoa(*port)

	log.Printf("Server will accept connections on %s...", address)

	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Panicln(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go handleRequest(conn, *delim)
	}
}

func Calculate(accumulator int64, cmd string, arg int64) int64{
	switch cmd {
	case "ADD":
		accumulator += arg
	case "SUB":
		accumulator -= arg
	case "MUL":
		accumulator *= arg
	case "SET":
		accumulator = arg
	default:
		log.Printf("Unknown command: %s", cmd)
	}
	return accumulator
}

func handleRequest(conn net.Conn, delim string) {
	log.Println("Accepted new connection.")
	defer conn.Close()
	defer log.Println("Closed connection.")

	reader := bufio.NewReader(conn)
	for {
		accumulator := int64(0)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() != "EOF" {
					log.Printf("Read error: %v", err)
				}
				return
			}

			line = strings.TrimRight(line, delim)
			if line == "" {
				// End of sequence, send result
				fmt.Fprintf(conn, "%d\r\n", accumulator)
				break
			}

			parts := strings.SplitN(line, " ", 2) // seperate by space
			if len(parts) != 2 {
				log.Panicln("Usage: command_name('ADD', 'SUB', 'MUL', 'SET') argument_value")
			}
									//  argument   base   bit_size
			arg, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				log.Printf("Invalid argument: %s", parts[1])
				continue
			}

			accumulator = Calculate(accumulator, parts[0], arg)
		}
	}
}
