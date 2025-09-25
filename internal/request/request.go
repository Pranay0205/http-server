package request

import (
	"bytes"
	"errors"
	"fmt"
	"http-server/internal/headers"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       State
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type State int

const (
	requestStateInitialized State = iota
	requestStateDone
	requestStateParsingHeaders
)

const crlf = "\r\n"

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)

	readToIndex := 0

	req := Request{}
	req.Headers = headers.NewHeaders()
	req.state = requestStateInitialized

	for req.state != requestStateDone {

		if readToIndex >= len(buf) {
			new_buf := make([]byte, len(buf)*2)
			copy(new_buf, buf[:readToIndex])
			buf = new_buf
		}

		readBytes, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.state = requestStateDone
				break
			}
			return &req, err
		}

		readToIndex += readBytes

		totalBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			req.state = requestStateDone
			return &req, err
		}

		if totalBytesParsed > 0 {
			copy(buf, buf[totalBytesParsed:])
			readToIndex -= totalBytesParsed
		}

	}

	if req.state != requestStateDone {
		return nil, fmt.Errorf("incomplete request: no valid request line found")
	}

	return &req, nil
}

func (r *Request) parse(data []byte) (int, error) {

	switch r.state {

	case requestStateInitialized:
		requestLine, bytesRead, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if bytesRead == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.state = requestStateParsingHeaders

		return bytesRead, nil
	case requestStateDone:

		return 0, fmt.Errorf("invalid state of the request: tryin to read data in done state")

	case requestStateParsingHeaders:

		h := r.Headers

		totalBytesParsed, done, err := h.Parse(data)

		if err != nil {
			return 0, err
		}

		if done {
			r.state = requestStateDone
		}

		return totalBytesParsed, nil

	default:

		return 0, fmt.Errorf("unknown state of the request: %d", r.state)

	}
}

func parseRequestLine(content []byte) (*RequestLine, int, error) {

	idx := bytes.Index(content, []byte(crlf))
	if idx == -1 {
		return &RequestLine{}, 0, nil
	}

	line := string(content[:idx])

	reqLine, err := requestLineExtractor(line)

	if err != nil {
		return nil, 0, err
	}

	return reqLine, idx + len(crlf), nil

}

func requestLineExtractor(line string) (*RequestLine, error) {
	if len(line) == 0 {
		return nil, fmt.Errorf("invalid request format: no request line found")
	}

	requestLine := line
	parts := strings.Split(requestLine, " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line: expected 3 parts, got %d", len(parts))
	}

	httpMethod := parts[0]
	if httpMethod != strings.ToUpper(httpMethod) {
		return nil, fmt.Errorf("invalid HTTP method: must be uppercase")
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("invalid HTTP version: missing version parts %s", parts[2])
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("invalid HTTP version: unrecognized HTTP-version %s", httpPart)
	}

	httpVersion := versionParts[1]
	if httpVersion != "1.1" {
		return nil, fmt.Errorf("invalid HTTP version: unrecognized HTTP-version %s", httpVersion)
	}

	return &RequestLine{
		Method:        httpMethod,
		RequestTarget: requestTarget,
		HttpVersion:   httpVersion,
	}, nil
}
