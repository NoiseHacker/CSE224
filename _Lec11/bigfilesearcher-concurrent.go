package main

import (
	"fmt"
	"time"
)

var filesizes = [5]int{3, 1, 5, 4, 2}
var results = [5]string{"ancient", "mariner", "bridegroom", "harbour", "albatross"}

func searchInBigFile(resultchannel chan<- string, filenum int) string {
	for i := 0; i < filesizes[filenum]; i++ {
		fmt.Printf("Searching file %d...", filenum)

		time.Sleep(1 * time.Second)

		fmt.Printf("%3.0f%% complete\n", float64(i+1)/float64(filesizes[filenum])*100)
	}

	resultchannel <- results[filenum]

	return results[filenum]
}

func main() {
	resultchannel := make(chan string)

	for i := 0; i < 5; i++ {
		go searchInBigFile(resultchannel, i)
	}

	localresults := make([]string, 5)
	for i := 0; i < 5; i++ {
		localresults[i] = <-resultchannel
	}

	fmt.Printf("Found these results: %v\n", localresults)
}
