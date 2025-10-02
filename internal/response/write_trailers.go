package response

import (
	"fmt"
	"http-server/internal/headers"
	"log"
	"strings"
)

func (w *Writer) WriteTrailers(h headers.Headers) error {

	if w.writerState != stateBodyWritten {
		return fmt.Errorf("trailer already written or invalid state")
	}

	trailerValues := h.Get("Trailer")

	trailerNames := strings.Split(trailerValues, ",")

	for _, trailer := range trailerNames {
		key := strings.TrimSpace(trailer)
		lowerKey := strings.ToLower(key)

		v, ok := h[lowerKey]

		if !ok {

			continue
		}

		trailerline := key + ": " + v + crlf
		log.Printf("Writing trailer line: %q", trailerline)
		w.Body = append(w.Body, trailerline...)
	}

	w.Body = append(w.Body, crlf...)

	w.SetState(stateTrailerWritten)
	return nil
}
