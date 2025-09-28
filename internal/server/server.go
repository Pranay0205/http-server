package server

import (
	"bytes"
	"fmt"
	"http-server/internal/request"
	"http-server/internal/response"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type HandlerError struct {
	StatusCode   response.StatusCode
	ErrorMessage error
}

type HandlerFunc func(w io.Writer, req *request.Request) *HandlerError

type Server struct {
	listener net.Listener
	enabled  *atomic.Bool
	handler  HandlerFunc
}

func Serve(port int, handlerFunc HandlerFunc) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: l,
		enabled:  &atomic.Bool{},
		handler:  handlerFunc,
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
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		writehandleError(conn, HandlerError{StatusCode: response.StatusBadRequest, ErrorMessage: err})
		return
	}

	var wb = make([]byte, 0, 10485760)
	responseBody := bytes.NewBuffer(wb)

	handlerErr := s.handler(responseBody, req)

	if handlerErr != nil {

		writehandleError(conn, *handlerErr)
		return

	} else {

		err = response.WriteStatusLine(conn, response.StatusSuccess)
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}

		headers := response.GetDefaultHeaders(responseBody.Len())

		err = response.WriteHeaders(conn, headers)
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}

		_, err = conn.Write(responseBody.Bytes())
		if err != nil {
			log.Printf("Error writing body: %v", err)
		}
	}

}

func writehandleError(w io.Writer, handlerErr HandlerError) {

	err := response.WriteStatusLine(w, handlerErr.StatusCode)
	if err != nil {
		log.Printf("Error writing status line: %v", err)
		return
	}

	errMessage := handlerErr.ErrorMessage.Error()
	headers := response.GetDefaultHeaders(len(errMessage))

	err = response.WriteHeaders(w, headers)
	if err != nil {
		log.Printf("Error writing the headers: %v", err)
		return
	}

	_, err = w.Write([]byte(errMessage))

	if err != nil {
		log.Printf("Error writing to the writer: %v", err)
		return
	}
}
