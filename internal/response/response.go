package response

import (
	"fmt"
	"http-server/internal/headers"
	"log"
)

type StatusCode int
type State int

const (
	StatusSuccess       StatusCode = 200
	StatusBadRequest    StatusCode = 400
	StatusNotFound      StatusCode = 404
	StatusInternalError StatusCode = 500
)

const (
	stateInitialized State = iota
	stateStatusWritten
	stateHeadersWritten
	stateBodyWritten
)

type Writer struct {
	Headers     headers.Headers
	StatusCode  StatusCode
	Body        []byte
	writerState State
}

const crlf = "\r\n"

var statusLines = map[StatusCode][]byte{
	StatusSuccess:       []byte("HTTP/1.1 200 OK"),
	StatusBadRequest:    []byte("HTTP/1.1 400 Bad Request"),
	StatusNotFound:      []byte("HTTP/1.1 404 Not Found"),
	StatusInternalError: []byte("HTTP/1.1 500 Internal Server Error"),
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if w.writerState != stateInitialized {
		return fmt.Errorf("status line already written or invalid state")
	}

	statusLine, exists := statusLines[statusCode]

	if !exists {
		return fmt.Errorf("unsupported status code: %d", statusCode)
	}

	statusLineCopy := make([]byte, len(statusLine))
	copy(statusLineCopy, statusLine)
	statusLineCopy = append(statusLineCopy, crlf...)

	log.Printf("Writing Statusline to response: %s", statusLineCopy)
	w.Body = append(w.Body, statusLineCopy...)

	defer w.SetState(stateStatusWritten)

	return nil
}

func (w *Writer) Header() headers.Headers {
	return w.Headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	if w.writerState != stateStatusWritten {
		return fmt.Errorf("must write status line before headers")
	}

	for key, value := range headers {
		headerLine := key + ": " + value + crlf
		log.Printf("key: %s, value: %s", key, value)
		log.Printf("Writing Headerline to response: %s", headerLine)
		w.Body = append(w.Body, headerLine...)
	}
	w.Body = append(w.Body, crlf...)

	defer w.SetState(stateHeadersWritten)

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	if w.writerState != stateHeadersWritten {
		return 0, fmt.Errorf("must write headers before body")
	}

	log.Printf("Writing body to response: %s", p)
	w.Body = append(w.Body, p...)

	defer w.SetState(stateBodyWritten)

	return len(p), nil
}

func (w *Writer) GetDefaultHeaders(contentLen int) error {
	if w.writerState == stateInitialized {
		return fmt.Errorf("must write status line before headers")
	}
	w.Headers.Add("Content-Length", fmt.Sprintf("%d", contentLen))

	w.Headers.Add("Connection", "close")

	return nil
}

func NewWriter() *Writer {
	return &Writer{
		Headers:     headers.NewHeaders(),
		writerState: stateInitialized,
	}
}

func (w *Writer) SetState(s State) {
	w.writerState = s
}
