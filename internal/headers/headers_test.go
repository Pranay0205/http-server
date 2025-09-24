package headers

import (
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
	assert.Equal(t, "localhost:42069", headers["Host"])
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
		assert.Equal(t, "", headers["Host"])
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
		assert.Equal(t, "localhost", headers["Host"])
		assert.Equal(t, 16, n)
		assert.False(t, done)
	})

	t.Run("Header with multiple spaces after colon", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host:    localhost    \r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "localhost", headers["Host"])
		assert.Equal(t, 24, n)
		assert.False(t, done)
	})

	t.Run("Header with complex value", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Accept: text/html, application/json, */*\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "text/html, application/json, */*", headers["Accept"])
		assert.Equal(t, 42, n)
		assert.False(t, done)
	})

}
