package response

import (
	"fmt"
	"http-server/internal/headers"
	"io"
	"log"
)

type StatusCode int

const crlf = "\r\n"

const (
	StatusSuccess       StatusCode = 200
	StatusBadRequest    StatusCode = 400
	StatusNotFound      StatusCode = 404
	StatusInternalError StatusCode = 500
)

var statusLines = map[StatusCode][]byte{
	StatusSuccess:       []byte("HTTP/1.1 200 OK"),
	StatusBadRequest:    []byte("HTTP/1.1 400 Bad Request"),
	StatusNotFound:      []byte("HTTP/1.1 404 Not Found"),
	StatusInternalError: []byte("HTTP/1.1 500 Internal Server Error"),
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine, exists := statusLines[statusCode]

	if !exists {
		return fmt.Errorf("unsupported status code: %d", statusCode)
	}

	statusLineCopy := make([]byte, len(statusLine))
	copy(statusLineCopy, statusLine)
	statusLineCopy = append(statusLineCopy, crlf...)

	_, err := w.Write(statusLineCopy)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := make(headers.Headers)

	header.Add("Content-Length", fmt.Sprintf("%d", contentLen))

	header.Add("Connection", "close")

	header.Add("Content-Type", "text/plain")

	return header
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		headerLine := key + ": " + value + crlf
		log.Printf("Headerline : %s", headerLine)
		_, err := w.Write([]byte(headerLine))
		if err != nil {

			log.Printf("Error From Write Headers: %s", err)
			return fmt.Errorf("unable to write headers to response: %w", err)
		}
	}

	_, err := w.Write([]byte(crlf))
	if err != nil {
		return fmt.Errorf("unable to write header separator: %w", err)
	}

	return nil
}
