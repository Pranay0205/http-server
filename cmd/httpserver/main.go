package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"http-server/internal/request"
	"http-server/internal/response"
	"http-server/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println()
	log.Println("Server gracefully stopped")
}

func handleClientError(w *response.Writer) {
	body := []byte("<html><head><title>400 Bad Request</title></head><body><h1>400 Bad Request</h1><p>The server could not understand the request due to invalid syntax.</p></body></html>")
	w.StatusCode = response.StatusBadRequest
	w.WriteStatusLine(w.StatusCode)
	w.GetDefaultHeaders(len(body))
	w.WriteHeaders(w.Headers)
	w.WriteBody(body)
}

func handleServerError(w *response.Writer) {
	body := []byte("<html><head><title>500 Internal Server Error</title></head><body><h1>500 Internal Server Error</h1><p>The server encountered an unexpected condition that prevented it from fulfilling the request.</p></body></html>")
	w.StatusCode = response.StatusInternalError
	w.WriteStatusLine(w.StatusCode)
	w.GetDefaultHeaders(len(body))
	w.WriteHeaders(w.Headers)
	w.WriteBody(body)
}

func handleProxyError(w *response.Writer, err error, newURL string) {
	errorMessage := fmt.Sprintf("Error while proxying request to %s: %s", newURL, err)
	body := []byte(fmt.Sprintf("<html><head><title>500 Internal Server Error</title></head><body><h1>500 Internal Server Error</h1><p>%s</p></body></html>", errorMessage))

	w.StatusCode = response.StatusInternalError
	w.WriteStatusLine(w.StatusCode)
	w.GetDefaultHeaders(len(body))
	w.WriteHeaders(w.Headers)
	w.WriteChunkedBody(body)
	log.Printf("couldn't reach the site: %s\n", newURL)
}

func streamProxyResponse(w *response.Writer, res *http.Response) {
	w.StatusCode = response.StatusSuccess
	w.WriteStatusLine(w.StatusCode)
	w.GetDefaultHeaders(int(res.ContentLength))
	w.Header().Set("Transfer-Encoding", "chunked")
	err := w.Header().Announce("X-Content-SHA256, X-Content-Length")
	if err != nil {
		log.Printf("Error announcing trailers: %v", err)
	}

	w.Header().Delete("Content-Length")
	w.Header().Override("Content-Type", "application/json")
	w.WriteHeaders(w.Headers)

	buf := make([]byte, 1024)
	responseBody := ""
	for {
		n, err := res.Body.Read(buf)
		log.Printf("read chunks = Hex: %X, Dec: %d bytes\n", n, n)
		if n > 0 {
			log.Println("encoded string: ", string(buf[:n]))
			_, err := w.WriteChunkedBody(buf[:n])
			if err != nil {
				log.Println("unable to write the body: ", err)
				break
			}
			responseBody += string(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error reading response body:", err)
			break
		}
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		log.Println("Error writing chunked body done:", err)
	}

	sha256 := fmt.Sprintf("%x", sha256.Sum256([]byte(responseBody)))

	err = w.Header().SetTrailer("X-Content-SHA256", sha256)
	if err != nil {
		log.Printf("Error setting SHA256 trailer: %v", err)
	}

	err = w.Header().SetTrailer("X-Content-Length", fmt.Sprintf("%d", len(responseBody)))
	if err != nil {
		log.Printf("Error setting length trailer: %v", err)
	}

	w.WriteTrailers(w.Headers)
}

func getProxyURL(requestTarget string) string {
	target := strings.TrimPrefix(requestTarget, "/")
	newURL := strings.Replace(target, "httpbin", "httpbin.org", 1)
	return "https://" + newURL
}

func handler(w *response.Writer, req *request.Request) {
	w.Header().Override("Content-Type", "text/html")

	if req.RequestLine.RequestTarget == "/yourproblem" {
		handleClientError(w)
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		handleServerError(w)
		return
	}

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyhandler(w, *req)
		return
	} else {

		body := []byte("<html><head><title>200 OK</title></head><body><h1>200 OK</h1><p>The request has succeeded.</p></body></html>")
		w.StatusCode = response.StatusSuccess
		w.WriteStatusLine(w.StatusCode)
		w.GetDefaultHeaders(len(body))
		w.WriteHeaders(w.Headers)
		w.WriteBody(body)
	}
}

func proxyhandler(w *response.Writer, req request.Request) {
	newURL := ""
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		newURL = getProxyURL(req.RequestLine.RequestTarget)
		log.Printf("Proxying to URL: %s\n", newURL)
	}
	res, err := http.Get(newURL)
	if err != nil {
		handleProxyError(w, err, newURL)
		return
	}
	streamProxyResponse(w, res)
}
