package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	port := flag.Int("port", 3333, "Port to connect to")
	host := flag.String("host", "127.0.0.1", "Host or IP to connect to")
	delim := flag.String("delimiter", "/", "Delimiter used to separate names")
	flag.Parse()

	c, err := net.Dial("tcp", *host+":"+strconv.Itoa(*port))

	if err != nil {
		log.Panicln(err)
	}

	log.Println("Connected to '"+*host+"' on port", strconv.Itoa(*port))
	defer c.Close()

	remaining := ""
	buf := make([]byte, 10)

	for {
		for strings.Contains(remaining, *delim) {
			idx := strings.Index(remaining, *delim)
			fmt.Println(remaining[:idx])
			remaining = remaining[idx+1:] //pass the delimiter
		}

		size, err := c.Read(buf)
		if err != nil {
			break
		}
		data := buf[:size]
		remaining = remaining + string(data)
	}

	log.Println("Program ended with remaining = " + remaining)
}
