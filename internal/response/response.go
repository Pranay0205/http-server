package response

import (
	"http-server/internal/headers"
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
	stateTrailerWritten
)

type Writer struct {
	Headers     headers.Headers
	StatusCode  StatusCode
	Body        []byte
	writerState State
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
