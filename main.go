package main

import (
	"fmt"
	"log"
	"os"
)

const inputFilePath = "message.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
		return
	}
	defer file.Close()

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	for {
		buffer := make([]byte, 8)

		readBytes, err := file.Read(buffer)

		if readBytes == 0 {
			fmt.Printf("Unable to read the file: %s", err)
			return
		}

		fmt.Printf("read: %s\n", buffer)
	}

}
