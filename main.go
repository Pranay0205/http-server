package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {

	channel := make(chan string)

	go func() {
		currentLine := ""
		defer f.Close()
		for {
			buffer := make([]byte, 8)

			readBytes, err := f.Read(buffer)

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

const port = ":42069"

func main() {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("❌ error listening for TCP traffic: %s\n", err.Error())
		return
	}

	fmt.Printf("👂 TCP Listener on Port:%s\n", listener.Addr())
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("❌ %s", err)
		}

		fmt.Printf("🔗 Connection established :%s\n", conn.RemoteAddr())

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Printf("✍️  %s\n", line)
		}
		log.Println("⛓️‍💥 Connection to ", conn.RemoteAddr(), "closed")
	}

}
