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

// Parses the header line
func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	key, value, n, done, err := parseHeaderLine(data)
	if err != nil {
		return 0, done, err
	}

	if h.Has(key) {
		h.Add(key, value)
	} else {
		h.Set(key, value)
	}

	return n, done, nil
}

// Parses header line and validates the key
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

// Extracts key, value from the headerline and validates formatting
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

// Set adds or updates a header (case-insensitive)
func (h Headers) Set(key, value string) {
	h[strings.ToLower(key)] = value
}

// Add appends a value to existing header or creates new one (case-insensitive)
func (h Headers) Add(key, value string) {
	lowerKey := strings.ToLower(key)
	if existing, ok := h[lowerKey]; ok {
		h[lowerKey] = existing + ", " + value
	} else {
		h[lowerKey] = value
	}
}

// Get retrieves header value (case-insensitive)
func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

// Has checks if header exists (case-insensitive)
func (h Headers) Has(key string) bool {
	_, ok := h[strings.ToLower(key)]
	return ok
}

// Delete removes header (case-insensitive)
func (h Headers) Delete(key string) {
	delete(h, strings.ToLower(key))
}
