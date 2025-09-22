package main

import (
	"fmt"
	"http-server/internal/request"
	"log"
	"net"
)

const port = ":42069"

func main() {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error listening for TCP traffic: %s\n", err.Error())
		return
	}

	fmt.Printf("👂 TCP Listener on Port:%s\n", listener.Addr())
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("%s", err)
		}

		reqLine, err := request.RequestFromReader(conn)

		if err != nil {
			log.Fatalf("%s", err)
		}

		fmt.Printf("Request line:\n - Method: %s\n - Target: %s\n - Version: %s\n", reqLine.RequestLine.Method, reqLine.RequestLine.RequestTarget, reqLine.RequestLine.HttpVersion)

		fmt.Printf("Connection established :%s\n", conn.RemoteAddr())

		log.Println("Connection to ", conn.RemoteAddr(), "closed")
	}

}
