package request

import (
	"bytes"
	"errors"
	"fmt"
	"http-server/internal/headers"
	"io"
	"strconv"
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
	requestStateParsingBody
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
		// Parsing residual part in the buffer before reading it from the connection or reader.
		// Prevents the blocking of the reader as it waits forever to read
		if readToIndex > 0 {
			for {
				totalBytesParsed, err := req.parse(buf[:readToIndex])
				if err != nil {
					return nil, err
				}

				if totalBytesParsed == 0 {
					break
				}

				if totalBytesParsed > 0 {
					copy(buf, buf[totalBytesParsed:])
					readToIndex -= totalBytesParsed
				}

				if req.state == requestStateDone {
					return &req, nil
				}
			}
		}

		// double checking if we are done with the request
		if req.state == requestStateDone {
			break
		}

		// If buffer is full doubling it
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf[:readToIndex])
			buf = newBuf
		}

		readBytes, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {

				readToIndex += readBytes

				for readToIndex > 0 {
					totalBytesParsed, parseErr := req.parse(buf[:readToIndex])
					if parseErr != nil {
						return nil, parseErr
					}

					if totalBytesParsed == 0 {
						break
					}

					copy(buf, buf[totalBytesParsed:])
					readToIndex -= totalBytesParsed
				}

				if req.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request: in state %d", req.state)
				}
				break
			}
			return nil, err
		}

		readToIndex += readBytes
	}

	return &req, nil
}

func (req *Request) parse(data []byte) (int, error) {

	switch req.state {

	case requestStateInitialized:
		requestLine, bytesRead, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if bytesRead == 0 {
			return 0, nil
		}

		req.RequestLine = *requestLine
		req.state = requestStateParsingHeaders

		return bytesRead, nil
	case requestStateDone:

		return 0, fmt.Errorf("invalid state of the request: tryin to read data in done state")

	case requestStateParsingHeaders:

		h := req.Headers

		totalBytesParsed, done, err := h.Parse(data)

		if err != nil {
			return 0, err
		}

		if done {
			req.state = requestStateParsingBody
		}

		return totalBytesParsed, nil

	case requestStateParsingBody:
		contentValue, exists := req.Headers["content-length"]

		if !exists {
			if len(data) > 0 && req.RequestLine.Method != "GET" {
				return 0, fmt.Errorf("invalid request: Content-Length header missing for non-empty body :%q", data)
			}
			req.state = requestStateDone
			return len(data), nil
		}

		req.Body = append(req.Body, data...)

		contentLength, err := strconv.Atoi(contentValue)

		if err != nil {
			return 0, err
		}

		if len(req.Body) > contentLength {
			return 0, fmt.Errorf("invalid request: body length exceeds Content-Length header")
		}

		if len(req.Body) == contentLength {
			req.state = requestStateDone
		}

		return len(data), nil

	default:

		return 0, fmt.Errorf("unknown state of the request: %d", req.state)

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
