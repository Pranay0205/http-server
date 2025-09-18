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
	if err != nil {
		return nil, fmt.Errorf("invalid request: unable to read content: %w", err)
	}
	if len(content) == 0 {
		return nil, fmt.Errorf("invalid request: empty content")
	}

	req, err := parseRequestLine(content)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func parseRequestLine(content []byte) (*RequestLine, error) {
	if len(content) == 0 {
		return nil, fmt.Errorf("invalid request: empty content")
	}

	lines := strings.Split(string(content), "\r\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("invalid request format: no lines found")
	}

	requestLine := lines[0]
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
