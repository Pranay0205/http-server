package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(file io.ReadCloser) <-chan string {

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	channel := make(chan string)

	go func() {
		currentLine := ""
		defer file.Close()
		for {
			buffer := make([]byte, 8)

			readBytes, err := file.Read(buffer)

			if err != nil {
				if currentLine != "" {
					channel <- currentLine
				}
				close(channel)
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("End of the line:%s", err.Error())
				return
			}

			data := currentLine + string(buffer[:readBytes])
			parts := strings.Split(data, "\n")

			for i := 0; i < len(parts)-1; i++ {
				channel <- parts[i]
			}

			currentLine = parts[len(parts)-1]

		}

	}()

	return channel
}

const inputFilePath = "message.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
		return
	}

	channelVariable := getLinesChannel(file)

	for line := range channelVariable {
		fmt.Printf("read: %s\n", line)
	}

}
