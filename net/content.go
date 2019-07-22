package net

import (
	"net/http"
	"strings"
)

// HasBinaryContent check that response or request has "binary" (non-text) content.
func HasBinaryContent(header http.Header, bytes []byte) bool {
	if contentType := header.Get("Content-Type"); contentType == "multipart/form-data" {
		return true
	}

	if contentType := header.Get("Content-Disposition"); contentType == "attachment" {
		return true
	}

	mimeType := http.DetectContentType(bytes)

	return false == strings.HasPrefix(mimeType, "text")
}
