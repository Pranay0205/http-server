package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const addr = "localhost:42069"

func main() {
	raddr, err := net.ResolveUDPAddr("udp", addr)

	if err != nil {
		log.Fatalf("😭 error listening for UDP traffic:", err)
		return
	}

	fmt.Printf("✅ UDP Addr is ready on %s\n\n", raddr)

	conn, err := net.DialUDP(raddr.Network(), nil, raddr)
	if err != nil {
		log.Fatalf("😭 error dialing UDP: %s\n", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")

		text, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalf("😭 Error reading the input: %s\n", err)
			return
		}

		_, err = conn.Write([]byte(text))
		if err != nil {
			log.Fatalf("😭 Error sending the message: %s\n", err)
		}
	}

}
