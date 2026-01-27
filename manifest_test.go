package asset

import (
	"bytes"
	"testing"
)

func TestWebpackParser(t *testing.T) {
	a := NewAssetMapper()
	a.PublicPath = ""

	var buff bytes.Buffer
	buff.WriteString(`
{
  "bundle.css": "/bundle/bundle.css",
  "bundle.js": "/bundle/bundle.06150a645f5d9e7e714c.js",
  "runtime.js": "/bundle/runtime.4243ae4e630918a98c56.js"
}`)

	table := []struct {
		in  string
		out string
	}{
		{"bundle.css", "/bundle/bundle.css"},
		{"bundle.js", "/bundle/bundle.06150a645f5d9e7e714c.js"},
	}

	readWebpackManifest(&buff, a)
	var result string

	for _, tt := range table {
		result = a.Get(tt.in)
		if tt.out != result {
			t.Errorf("webpack manifest parse err:\n expected: %s\ngot: %s\n", tt.out, result)
		}
	}
}

func TestViteParser(t *testing.T) {
	a := NewAssetMapper()
	a.PublicPath = ""

	var buff bytes.Buffer
	buff.WriteString(`
{
  "src/app.js": {
    "file": "assets/app-CKgRTByK.js",
    "name": "app",
    "src": "src/app.js",
    "isEntry": true,
    "css": [
      "assets/app-o2N34dPp.css"
    ]
  },
  "src/placeholder.jpeg": {
    "file": "assets/placeholder-DXdl7YkJ.jpeg",
    "src": "src/placeholder.jpeg"
  }
}`)

	table := []struct {
		in  string
		out string
	}{
		{"src/app.js", "assets/app-CKgRTByK.js"},
		{"src/placeholder.jpeg", "assets/placeholder-DXdl7YkJ.jpeg"},
	}

	readViteManifest(&buff, a)
	var result string

	for _, tt := range table {
		result = a.Get(tt.in)
		if tt.out != result {
			t.Errorf("vite manifest parse err:\n expected: %s\ngot: %s\n", tt.out, result)
		}
	}
}

