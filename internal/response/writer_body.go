package response

import (
	"fmt"
	"log"
)

func (w *Writer) WriteBody(p []byte) (int, error) {

	if w.writerState != stateHeadersWritten {
		return 0, fmt.Errorf("must write headers before body")
	}

	log.Printf("Writing body to response: %s", p)
	w.Body = append(w.Body, p...)

	defer w.SetState(stateBodyWritten)

	return len(p), nil
}
