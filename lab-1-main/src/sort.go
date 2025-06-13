package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"slices"
	"bytes"
)

type Records struct {
	length   uint32
	key  [10]byte
	value []byte
}

/* Read a big-endian uint32 from a byte slice of length at least 4*/
func ReadBigEndianUint32(buffer []byte) uint32 {
	if len(buffer) < 4 {
		panic("buffer too short to read uint32")
	}
	return binary.BigEndian.Uint32(buffer[:])
}

func SplitRecords(data []byte) []Records {
	var records []Records
	offset := 0

	for offset < len(data) {
		// Split length
		lengths := ReadBigEndianUint32(data[offset : offset+4])
		offset += 4

		// validate the length
		if lengths < 10 {
			fmt.Printf("invalid length")
		}

		//validate the completeness of record 
		endOfRecord := offset + int(lengths)
		if endOfRecord > len(data) {
			fmt.Printf("incomplete record")
		}

		// Split key 
		var keys [10]byte
		copy(keys[:], data[offset:offset+10])

		// Split value
		values := data[offset+10 : endOfRecord]

		records = append(records, Records{length: lengths, key: keys, value: values})
		offset = endOfRecord
	}

	return records
}

// Write a big-endian uint32 to a byte slice of length at least 4
func WriteBigEndianUint32(buffer []byte, num uint32) {
	if len(buffer) < 4 {
		panic("buffer too short to write uint32")
	}
	binary.BigEndian.PutUint32(buffer, num)
}

// Sort Custom Record types
func CompareRecords(a, b Records) int {
	return bytes.Compare(a.key[:], b.key[:])
}

func main() {
	log.Printf("")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) != 3 {
		log.Fatalf("Usage: %v inputfile outputfile\n", os.Args[0])
	}

	log.Printf("Sorting %s to %s\n", os.Args[1], os.Args[2])

	//read bytes from a file
	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	records := SplitRecords(input)

	//sort records
	slices.SortFunc(records, CompareRecords)

	// Create a file to write
	output, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Println(err)
	}
	defer output.Close()

	for _, rec := range records {
		// length contains the sum of the lengths of the key and value fields
		// length := uint32(len(rec.key)) + uint32(len(rec.value))
		buffer := make([]byte, 4)
		WriteBigEndianUint32(buffer, rec.length)

		// write length to output file
		_, err := output.Write(buffer)
		if err != nil {
			fmt.Println(err)
		}

		// write key to output file
		_, err = output.Write(rec.key[:]) 
		if err != nil {
			fmt.Println(err)
		}

		// write value to output file
		_, err = output.Write(rec.value)
		if err != nil {
			fmt.Println(err)
		}
	}

}
