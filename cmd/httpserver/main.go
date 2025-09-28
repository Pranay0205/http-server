package main

import (
	"errors"
	"fmt"
	"http-server/internal/request"
	"http-server/internal/response"
	"http-server/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {

	handler := func(w io.Writer, req *request.Request) *server.HandlerError {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			return &server.HandlerError{StatusCode: response.StatusBadRequest, ErrorMessage: errors.New("Your problem is not my problem")}
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			return &server.HandlerError{StatusCode: response.StatusInternalError, ErrorMessage: errors.New("Woopsie, my bad")}
		}

		w.Write([]byte("All good, frfr"))
		return nil
	}

	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	fmt.Println()
	log.Println("Server gracefully stopped")
}
