package server

import (
	"fmt"
	"http-server/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	enabled  *atomic.Bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: l,
		enabled:  &atomic.Bool{},
	}
	s.enabled.Store(true)

	go s.listen()
	return s, nil
}

func (s *Server) listen() {
	for s.enabled.Load() {
		conn, err := s.listener.Accept()
		if err != nil {

			if !s.enabled.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) Close() error {
	s.enabled.Store(false)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) handle(conn net.Conn) {
	err := response.WriteStatusLine(conn, response.StatusSuccess)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}

	headers := response.GetDefaultHeaders(0)

	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
	defer conn.Close()
}
