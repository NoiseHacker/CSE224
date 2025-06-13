package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	pb "numlog/numlog"
	"os"

	"google.golang.org/protobuf/proto"
)

func main() {
	logfp, err := os.Open("log.db")
	if err != nil {
		log.Panic(err)
	}
	defer logfp.Close()

	for {
		loglenbytes := make([]byte, 4)
		_, err := io.ReadFull(logfp, loglenbytes)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		loglen := binary.LittleEndian.Uint32(loglenbytes)

		logbytes := make([]byte, loglen)
		_, err = io.ReadFull(logfp, logbytes)
		if err != nil {
			log.Panic(err)
		}

		logentry := &pb.NumberEntry{}

		if err := proto.Unmarshal(logbytes, logentry); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("We received a number at time %v\n", logentry.Ts)

	}
}
