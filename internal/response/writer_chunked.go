package response

import (
	"fmt"
	"log"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != stateHeadersWritten {
		return 0, fmt.Errorf("must write headers before body")
	}

	before := len(w.Body)

	log.Printf("Writing body to response: %s\n", p)
	chunk := fmt.Sprintf("%X\r\n", len(p))
	w.Body = append(w.Body, chunk...)
	w.Body = append(w.Body, p...)
	w.Body = append(w.Body, []byte(crlf)...)

	after := len(w.Body)

	return after - before, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {

	if w.writerState != stateHeadersWritten {
		return 0, fmt.Errorf("must write headers before body")
	}
	body := []byte("0\r\n")
	log.Printf("end of the file reached!\n")
	w.Body = append(w.Body, body...)

	defer w.SetState(stateBodyWritten)

	return len(body), nil
}
