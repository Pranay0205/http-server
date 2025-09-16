package main

import (
	"fmt"
	"log"
	"os"
	"strings"
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

	currentLine := ""
	for {
		buffer := make([]byte, 8)

		readBytes, err := file.Read(buffer)

		if readBytes == 0 {
			fmt.Printf("Unable to read the file: %s", err)
			break
		}

		data := currentLine + string(buffer[:readBytes])
		parts := strings.Split(data, "\n")

		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s\n", parts[i])
		}

		currentLine = parts[len(parts)-1]
	}

	if currentLine != "" {
		fmt.Printf("read: %s\n", currentLine)
	}
}
