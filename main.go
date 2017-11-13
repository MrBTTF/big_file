package main

import (
	"golang.org/x/exp/mmap"
	"fmt"
	"math"
	"encoding/binary"
	"bytes"
	"time"
)

var (
	memorySize = int(math.Pow10(8))
	numberToFind = 18643
)

func readChunk(file *mmap.ReaderAt, i int, bytesRead *int, result chan bool, quit chan bool) {
	buff := make([]byte, memorySize)

	n, err := file.ReadAt(buff,int64(i*memorySize))
	if err != nil {
		fmt.Println(err)
	}

	if findNumber(n,buff,quit){
		result <- true
		quit <- true
	} else {
		result <- false
	}


	*bytesRead += n
	if *bytesRead % (100000000) == 0 {
		fmt.Printf("Read %d megabytes\n", *bytesRead/1000000)
	}

}

func findNumber(n int,buff []byte,quit chan bool) bool {
	for j:=0; j < n; j+=8 {
		select{
		case <-quit:
			return false
		default:
			var num int64
			binary.Read(bytes.NewBuffer(buff[j:j+8]), binary.BigEndian, &num)
			if num == int64(numberToFind) {
				return true
			}
		}
	}
	return false
}

func main() {
	start := time.Now()

	file, err := mmap.Open("data")
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("File size: %d bytes\n",file.Len())

	chunksCount := file.Len() / memorySize
	if file.Len() % memorySize !=0 {
		chunksCount += 1
	}
	fmt.Printf("Chunks count: %d chunks\n",chunksCount)

	found := false
	bytesRead := 0
	result := make(chan bool)
	quit := make(chan bool)
	for i := 0; i < chunksCount; i++{
		go readChunk(file, i, &bytesRead,result,quit)
	}

	for b := range result {
		fmt.Println(b)
		if b{
			fmt.Printf("Number %d found\n",numberToFind)
			found = true
			break
		}

	}
	if !found{
		fmt.Printf("Number %d not found\n",numberToFind)
	}

	fmt.Printf("Time took %s\n", time.Since(start))
}
