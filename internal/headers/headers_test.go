package headers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Error/Edge Cases
	t.Run("Missing colon", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host localhost\r\n")
		n, done, err := headers.Parse(data)
		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})

	t.Run("Empty key", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte(": value\r\n")
		n, done, err := headers.Parse(data)
		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})

	t.Run("Empty value", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host:\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "", headers["host"])
		assert.Equal(t, 7, n)
		assert.False(t, done)
	})

	t.Run("No CRLF found", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})

	// Invalid Spacing Cases
	t.Run("Space before and after colon", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host : localhost\r\n")
		n, done, err := headers.Parse(data)
		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})

	t.Run("Multiple spaces before colon", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host  : localhost\r\n")
		n, done, err := headers.Parse(data)
		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})

	// Valid Edge Cases
	t.Run("Header with no space after colon", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host:localhost\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "localhost", headers["host"])
		assert.Equal(t, 16, n)
		assert.False(t, done)
	})

	t.Run("Header with multiple spaces after colon", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host:    localhost    \r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "localhost", headers["host"])
		assert.Equal(t, 24, n)
		assert.False(t, done)
	})

	t.Run("Header with complex value", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Accept: text/html, application/json, */*\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "text/html, application/json, */*", headers["accept"])
		assert.Equal(t, 42, n)
		assert.False(t, done)
	})

	t.Run("Header with uppercase key", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("HOST: localhost:4337\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "localhost:4337", headers["host"])
		assert.Empty(t, headers["HOST"])

		assert.Equal(t, len(data), n)
		assert.False(t, done)
	})

	t.Run("Header with invalid key", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("H©st: localhost:42069\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})

	t.Run("Set and Get with different cases", func(t *testing.T) {
		headers := make(Headers)
		headers.Set("Content-Type", "application/json")

		// All these should return the same value
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, "application/json", headers.Get("content-type"))
		assert.Equal(t, "application/json", headers.Get("CONTENT-TYPE"))
		assert.Equal(t, "application/json", headers.Get("Content-type"))
	})

	t.Run("Typical HTTP request headers", func(t *testing.T) {
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
	})

	t.Run("Two Accept headers concatenated", func(t *testing.T) {
		headers := NewHeaders()

		// First Accept header
		data1 := []byte("Accept: text/html\r\n")
		n1, done1, err1 := headers.Parse(data1)
		require.NoError(t, err1)
		assert.Equal(t, "text/html", headers.Get("accept"))
		fmt.Printf("This is the %s", headers.Get("accept"))
		assert.Equal(t, 19, n1)
		assert.False(t, done1)

		// Second Accept header (should concatenate with comma-space)
		data2 := []byte("Accept: application/json\r\n")
		n2, done2, err2 := headers.Parse(data2)
		require.NoError(t, err2)
		assert.Equal(t, "text/html, application/json", headers.Get("accept"))
		assert.Equal(t, 26, n2)
		assert.False(t, done2)
	})

}
