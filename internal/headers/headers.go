package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {

	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	key, value, n, done, err := parseHeaderLine(data)
	if err != nil {
		return 0, done, err
	}

	key = strings.ToLower(key)

	h[key] = value

	return n, done, nil
}

func parseHeaderLine(content []byte) (key, value string, n int, done bool, err error) {

	idx := bytes.Index(content, []byte(crlf))
	if idx == -1 {
		return "", "", 0, false, nil
	}

	if idx == 0 {
		return "", "", idx + len(crlf), true, nil
	}

	line := string(content[:idx])

	key, value, err = headerLineExtractor(line)

	if err != nil {
		return "", "", 0, false, err
	}

	invalidChars := regexp.MustCompile(`[^a-zA-Z0-9!#$%&'*\+\-\.^_` + "`" + `|~]+`)

	if invalidChars.MatchString(key) {
		return "", "", 0, false, fmt.Errorf("invalid header name")
	}

	return key, value, idx + len(crlf), false, nil

}

func headerLineExtractor(headerLine string) (key, value string, err error) {

	headerLine = strings.Trim(headerLine, crlf)

	parts := strings.SplitN(headerLine, ":", 2)

	if len(parts) == 1 {
		return "", "", fmt.Errorf("no split element found")
	}

	if parts[0] == "" {
		return "", "", fmt.Errorf("empty key not allowed")
	}

	if parts[0] != strings.TrimSpace(parts[0]) {
		return "", "", fmt.Errorf("invalid header format: additional space detected")
	}

	key, value = parts[0], strings.TrimSpace(parts[1])

	return key, value, nil
}
