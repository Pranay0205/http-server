package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	content, err := io.ReadAll(reader)
	if len(content) == 0 {
		return nil, fmt.Errorf("HTTP request is empty")
	}
	if len(content) == 0 || err != nil {
		return nil, fmt.Errorf("unable to read the request line: %w", err)
	}

	req, err := parseRequestLine(content)

	if err != nil {
		return nil, fmt.Errorf("unable to parse the request: %w", err)
	}

	return req, nil
}

func parseRequestLine(content []byte) (*Request, error) {

	if len(content) == 0 {
		return nil, fmt.Errorf("HTTP Message is empty")
	}

	lines := strings.Split(string(content), "\r\n")

	if len(lines) == 0 {
		return nil, fmt.Errorf("HTTP Message is not in correct format")
	}

	requestLine := lines[0]

	parts := strings.Split(requestLine, " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid HTTP Message")
	}

	if parts[0] != strings.ToUpper(parts[0]) {
		return nil, fmt.Errorf("HTTP method is not in correct format")
	}

	if !strings.Contains(parts[2], "HTTP/1.1") {
		return nil, fmt.Errorf("HTTP Message version unsupported")
	}

	req := Request{}

	req.RequestLine.Method = parts[0]
	req.RequestLine.HttpVersion = strings.Split(parts[2], "/")[1]
	req.RequestLine.RequestTarget = parts[1]

	return &req, nil
}
