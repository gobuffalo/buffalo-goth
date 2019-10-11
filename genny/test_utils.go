package genny

import (
	"strings"
)

/**
 * Source: https://www.programming-books.io/essential/go/normalize-newlines-1d3abcf6f17c4186bb9617fa14074e48
 */
func NormalizeNewlines(d string) string {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = strings.ReplaceAll(d, string([]byte{13, 10}), string([]byte{10}))
	// replace CF \r (mac) with LF \n (unix)
	d = strings.ReplaceAll(d, string([]byte{13}), string([]byte{10}))
	return d
}
