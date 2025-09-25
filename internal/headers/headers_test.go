package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Basic Valid Cases
func TestHeaders_ValidSingleHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestHeaders_HeaderWithNoSpaceAfterColon(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host:localhost\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost", headers["host"])
	assert.Equal(t, 16, n)
	assert.False(t, done)
}

func TestHeaders_HeaderWithMultipleSpacesAfterColon(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host:    localhost    \r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost", headers["host"])
	assert.Equal(t, 24, n)
	assert.False(t, done)
}

func TestHeaders_HeaderWithComplexValue(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Accept: text/html, application/json, */*\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "text/html, application/json, */*", headers["accept"])
	assert.Equal(t, 42, n)
	assert.False(t, done)
}

func TestHeaders_HeaderWithUppercaseKey(t *testing.T) {
	headers := NewHeaders()
	data := []byte("HOST: localhost:4337\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:4337", headers["host"])
	assert.Empty(t, headers["HOST"])
	assert.Equal(t, len(data), n)
	assert.False(t, done)
}

func TestHeaders_EmptyValue(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host:\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "", headers["host"])
	assert.Equal(t, 7, n)
	assert.False(t, done)
}

// Error Cases
func TestHeaders_InvalidSpacingHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestHeaders_MissingColon(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host localhost\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestHeaders_EmptyKey(t *testing.T) {
	headers := NewHeaders()
	data := []byte(": value\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestHeaders_SpaceBeforeAndAfterColon(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host : localhost\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestHeaders_MultipleSpacesBeforeColon(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host  : localhost\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestHeaders_InvalidKeyCharacters(t *testing.T) {
	headers := NewHeaders()
	data := []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

// Edge Cases
func TestHeaders_NoCRLFFound(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

// Header Methods Tests
func TestHeaders_SetAndGetWithDifferentCases(t *testing.T) {
	headers := make(Headers)
	headers.Set("Content-Type", "application/json")

	// All these should return the same value
	assert.Equal(t, "application/json", headers.Get("Content-Type"))
	assert.Equal(t, "application/json", headers.Get("content-type"))
	assert.Equal(t, "application/json", headers.Get("CONTENT-TYPE"))
	assert.Equal(t, "application/json", headers.Get("Content-type"))
}

// Multiple Headers Tests
func TestHeaders_TypicalHTTPRequestHeaders(t *testing.T) {
	headers := NewHeaders()

	testData := [][]byte{
		[]byte("Host: api.example.com\r\n"),
		[]byte("User-Agent: MyApp/1.0\r\n"),
		[]byte("Accept: application/json\r\n"),
		[]byte("Accept: text/html\r\n"),
		[]byte("Content-Type: application/x-www-form-urlencoded\r\n"),
		[]byte("Authorization: Bearer abc123\r\n"),
	}

	for _, data := range testData {
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.False(t, done)
		assert.Greater(t, n, 0)
	}

	// Verify all headers are present
	assert.Equal(t, "api.example.com", headers.Get("Host"))
	assert.Equal(t, "MyApp/1.0", headers.Get("User-Agent"))
	assert.Equal(t, "application/json, text/html", headers.Get("Accept"))
	assert.Equal(t, "application/x-www-form-urlencoded", headers.Get("Content-Type"))
	assert.Equal(t, "Bearer abc123", headers.Get("Authorization"))
}

func TestHeaders_AcceptHeadersConcatenated(t *testing.T) {
	headers := NewHeaders()

	// First Accept header
	data1 := []byte("Accept: text/html\r\n")
	n1, done1, err1 := headers.Parse(data1)
	require.NoError(t, err1)
	assert.Equal(t, "text/html", headers.Get("accept"))
	assert.Equal(t, 19, n1)
	assert.False(t, done1)

	// Second Accept header (should concatenate with comma-space)
	data2 := []byte("Accept: application/json\r\n")
	n2, done2, err2 := headers.Parse(data2)
	require.NoError(t, err2)
	assert.Equal(t, "text/html, application/json", headers.Get("accept"))
	assert.Equal(t, 26, n2)
	assert.False(t, done2)
}

// Table-driven tests for invalid spacing patterns
func TestHeaders_InvalidSpacingPatterns(t *testing.T) {
	testCases := []struct {
		name string
		data string
	}{
		{
			name: "space before colon",
			data: "Host : localhost\r\n",
		},
		{
			name: "multiple spaces before colon",
			data: "Host  : localhost\r\n",
		},
		{
			name: "leading whitespace",
			data: "   Host: localhost\r\n",
		},
		{
			name: "trailing whitespace on key",
			data: "Host   : localhost\r\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers := NewHeaders()
			data := []byte(tc.data)
			n, done, err := headers.Parse(data)
			require.Error(t, err)
			assert.Equal(t, 0, n)
			assert.False(t, done)
		})
	}
}

// Table-driven tests for valid header formats
func TestHeaders_ValidHeaderFormats(t *testing.T) {
	testCases := []struct {
		name          string
		data          string
		expectedKey   string
		expectedValue string
		expectedBytes int
	}{
		{
			name:          "basic header",
			data:          "Host: localhost\r\n",
			expectedKey:   "host",
			expectedValue: "localhost",
			expectedBytes: 17,
		},
		{
			name:          "no space after colon",
			data:          "Host:localhost\r\n",
			expectedKey:   "host",
			expectedValue: "localhost",
			expectedBytes: 16,
		},
		{
			name:          "multiple spaces after colon",
			data:          "Host:   localhost   \r\n",
			expectedKey:   "host",
			expectedValue: "localhost",
			expectedBytes: 22,
		},
		{
			name:          "empty value",
			data:          "Host:\r\n",
			expectedKey:   "host",
			expectedValue: "",
			expectedBytes: 7,
		},
		{
			name:          "uppercase key",
			data:          "HOST: localhost\r\n",
			expectedKey:   "host",
			expectedValue: "localhost",
			expectedBytes: 17,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers := NewHeaders()
			data := []byte(tc.data)
			n, done, err := headers.Parse(data)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedValue, headers[tc.expectedKey])
			assert.Equal(t, tc.expectedBytes, n)
			assert.False(t, done)
		})
	}
}
