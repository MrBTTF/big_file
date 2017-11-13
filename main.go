package main

import (
	"golang.org/x/exp/mmap"
	"fmt"
	"math"
	"encoding/binary"
	"bytes"
)

var (
	memorySize = int(math.Pow10(3))
	numberToFind = 18643
)

func main() {
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

	buff := make([]byte, memorySize)
	bytesRead := 0
	for i := 0; i < chunksCount; i++{
		n, err := file.ReadAt(buff,int64(i*memorySize))
		if err != nil {
			fmt.Println(err)
			continue
		}

		bytesRead += n
		if bytesRead %(memorySize*1000) == 0 {
			fmt.Printf("Read %d bytes\n", bytesRead)
		}

		for j:=0; j < n; j+=8 {
			var num int64
			binary.Read(bytes.NewBuffer(buff[j:j+8]), binary.BigEndian, &num)
			if num == int64(numberToFind) {
				fmt.Printf("Number %d found in chunk %d index %d\n",numberToFind,i, j/8)
				return
			}
		}

	}

	fmt.Printf("Number %d not found\n",numberToFind)


}
