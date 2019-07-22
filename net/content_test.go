package net

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestHasBinaryContent_JSON(t *testing.T) {
	body, err := json.Marshal(struct {
		Name  string
		Value int
	}{"A", 1})
	if err != nil {
		t.Error(err)
	}

	if HasBinaryContent(http.Header{}, body) {
		t.Error("Should be text type")
	}
}

func TestHasBinaryContent_HTML(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<body>

<h1>My First Heading</h1>

<p>My first paragraph.</p>

</body>
</html>
`

	if HasBinaryContent(http.Header{}, []byte(html)) {
		t.Error("Should be text type")
	}
}

func TestHasBinaryContent_ZIP(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := zip.NewWriter(buf)
	fw, err := writer.Create("File")
	if err != nil {
		t.Fatal(err)
	}
	fw.Write([]byte("Hello World!"))
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}

	if !HasBinaryContent(http.Header{}, buf.Bytes()) {
		t.Error("Should be non-text type")
	}
}
