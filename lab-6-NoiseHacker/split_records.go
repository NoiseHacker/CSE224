package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	var inputFile string
	var outputDir string
	var numNodes int

	flag.StringVar(&inputFile, "input", "", "Path to input record file (from gensort)")
	flag.StringVar(&outputDir, "outdir", "inputs", "Directory to write per-node input files")
	flag.IntVar(&numNodes, "nodes", 4, "Number of nodes")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("Must provide --input")
	}

	in, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	outs := make([]*os.File, numNodes)
	for i := 0; i < numNodes; i++ {
		path := fmt.Sprintf("%s/input_%d.dat", outputDir, i)
		f, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		outs[i] = f
		defer f.Close()
	}

	recCount := 0
	for {
		var lenBuf [4]byte
		if _, err := io.ReadFull(in, lenBuf[:]); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		recLen := binary.BigEndian.Uint32(lenBuf[:])
		data := make([]byte, recLen)
		if _, err := io.ReadFull(in, data); err != nil {
			log.Fatal(err)
		}

		nodeID := recCount % numNodes

		outs[nodeID].Write(lenBuf[:])
		outs[nodeID].Write(data)
		recCount++
	}

	fmt.Printf("Done. Split %d records across %d nodes.\n", recCount, numNodes)
}
