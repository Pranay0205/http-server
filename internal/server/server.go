package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	Listner net.Listener
	enabled *atomic.Bool
}

func Serve(port int) (*Server, error) {

	l, err := net.Listen("tcp", ":"+fmt.Sprint(port))

	if err != nil {
		return nil, err
	}

	enabled := &atomic.Bool{}
	enabled.Store(true)

	server := Server{Listner: l, enabled: enabled}

	go func(server *Server) {
		server.listen()
	}(&server)

	return &server, nil

}

func (s *Server) listen() {
	for s.enabled.Load() {
		conn, err := s.Listner.Accept()

		if err != nil {
			if !s.enabled.Load() {
				return
			}
			s.Close()
			log.Fatalf("Unable to establish a connection")
		}

		go s.handle(conn)

	}
}

func (s *Server) Close() error {
	err := s.Listner.Close()
	s.enabled.Store(false)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handle(conn net.Conn) {
	message := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!")

	_, err := conn.Write(message)
	if err != nil {
		log.Printf("Unable to write the connection")
	}

	io.Copy(conn, conn)
	conn.Close()
}
