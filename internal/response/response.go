package response

import (
	"fmt"
	"http-server/internal/headers"
	"io"
	"strings"
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
	statusLine = append(statusLine, crlf...)
	_, err := w.Write(statusLine)

	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := make(headers.Headers)

	header.Add("Content-Length", fmt.Sprintf("%d", contentLen))

	header.Add("Connection", "close")

	header.Add("Content-Type", "text/plain")

	return header
}

func WriteHeader(w io.Writer, headers headers.Headers) error {
	headerStr := strings.Join([]string{
		headers["Content-Length"],
		headers["Connection"],
		headers["Content-Type"],
	}, crlf)

	_, err := w.Write([]byte(headerStr))

	if err != nil {
		return fmt.Errorf("unable to write headers to response: %w", err)
	}

	return nil
}
