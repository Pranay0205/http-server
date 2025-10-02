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

	defer w.SetState(stateTrailerWritten)

	trailerValues := w.Header().Get("Trailer")
	if trailerValues == "" {
		return fmt.Errorf("no trailers announced")
	}

	trailerNames := strings.Split(trailerValues, ",")

	for _, trailer := range trailerNames {
		key := strings.TrimSpace(trailer)
		lowerKey := strings.ToLower(key)

		v, ok := h[lowerKey]
		if !ok {
			log.Printf("announced trailer %s not found in provided headers", key)
			continue
		}

		trailerline := key + ": " + v + crlf
		log.Printf("Writing trailer line: %q", trailerline)
		w.Body = append(w.Body, trailerline...)
	}

	w.Body = append(w.Body, crlf...)

	return nil
}
