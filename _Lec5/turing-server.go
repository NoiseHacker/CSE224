package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {		// name  value 	usage
	port := flag.Int("port", 3333, "Port to accept connections on.")
	host := flag.String("host", "127.0.0.1", "Host or IP to bind to")
	filename := flag.String("file", "turing.txt", "File with Turing award names")
	delim := flag.String("delimiter", "/", "Delimiter that separates winners")
	flag.Parse()

	winners := readAwardFile(*filename)

	l, err := net.Listen("tcp", *host+":"+strconv.Itoa(*port))
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Listening to connections at '"+*host+"' on port", strconv.Itoa(*port))
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go handleRequest(conn, winners, *delim)
	}
}

func readAwardFile(filename string) []string {
	log.Println("Reading from file " + filename)

	f, err := os.Open(filename)
	if err != nil {
		log.Panicln(err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Panicln(err)
		}
	}()

	winners := []string{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		winners = append(winners, s.Text())
	}
	err = s.Err()
	if err != nil {
		log.Panicln(err)
	}

	return winners
}

// Writing data with delimiter
func handleRequest(conn net.Conn, winners []string, delim string) {
	log.Println("Accepted new connection.")
	defer conn.Close()
	defer log.Println("Closed connection.")

	for i := 0; i < len(winners); i++ {
		log.Println("Sending ", winners[i])

		valbytes := []byte(winners[i] + delim)
		_, err := conn.Write(valbytes)
		if err != nil {
			return
		}
	}
}
