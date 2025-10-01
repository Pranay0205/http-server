package response

import (
	"fmt"
	"log"
)

var statusLines = map[StatusCode][]byte{
	StatusSuccess:       []byte("HTTP/1.1 200 OK"),
	StatusBadRequest:    []byte("HTTP/1.1 400 Bad Request"),
	StatusNotFound:      []byte("HTTP/1.1 404 Not Found"),
	StatusInternalError: []byte("HTTP/1.1 500 Internal Server Error"),
}

func GetStatusLine(statusCode StatusCode) ([]byte, error) {
	v, ok := statusLines[statusCode]
	if !ok {
		return []byte{}, fmt.Errorf("unsupported status code")
	}

	return v, nil
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if w.writerState != stateInitialized {
		return fmt.Errorf("status line already written or invalid state")
	}

	statusLine, err := GetStatusLine(statusCode)

	if err != nil {
		return fmt.Errorf("unsupported status code: %d", statusCode)
	}

	statusLineCopy := make([]byte, len(statusLine))
	copy(statusLineCopy, statusLine)
	statusLineCopy = append(statusLineCopy, crlf...)

	log.Printf("Writing Statusline to response: %s\n", statusLineCopy)
	w.Body = append(w.Body, statusLineCopy...)

	defer w.SetState(stateStatusWritten)

	return nil
}
