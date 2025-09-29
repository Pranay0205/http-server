package response

import (
	"fmt"
	"http-server/internal/headers"
	"log"
)

const crlf = "\r\n"

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

func (w *Writer) GetDefaultHeaders(contentLen int) error {
	if w.writerState == stateInitialized {
		return fmt.Errorf("must write status line before headers")
	}
	w.Headers.Add("Content-Length", fmt.Sprintf("%d", contentLen))

	w.Headers.Add("Connection", "close")

	return nil
}
