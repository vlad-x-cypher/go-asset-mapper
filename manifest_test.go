package asset

import (
	"bytes"
	"testing"
)

func prepare(publicPath string) (*AssetMapper, bytes.Buffer) {
	a := NewAssetMapper()
	a.PublicPath = publicPath

	var buff bytes.Buffer

	return a, buff
}

func TestWebpackParser(t *testing.T) {
	a, buff := prepare("/static/")
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
		{"bundle.css", "/static/bundle/bundle.css"},
		{"bundle.js", "/static/bundle/bundle.06150a645f5d9e7e714c.js"},
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
	a, buff := prepare("/static/")
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
		{"src/app.js", a.PublicPath + "assets/app-CKgRTByK.js"},
		{"src/placeholder.jpeg", a.PublicPath + "assets/placeholder-DXdl7YkJ.jpeg"},
		{"non-existent.css", "non-existent.css"},
		// test correct entries parse
		{"assets/app-o2N34dPp.css", a.PublicPath + "assets/app-o2N34dPp.css"},
	}

	readViteManifest(&buff, a)
	var result string

	for _, tt := range table {
		result = a.Get(tt.in)
		if tt.out != result {
			t.Errorf("vite manifest parse err:\n expected: %s\ngot: %s\n", tt.out, result)
		}
	}

	// test entry
	entries := a.CSSEntry("app")
	if len(entries) < 1 {
		t.Errorf("entry count incorrect\n expected: %d\n got: %d", 1, len(entries))
	}
	entryCSS := a.Get("assets/app-o2N34dPp.css")
	if entries[0] != entryCSS {
		t.Errorf("entry path incorrect\n expected: %s\n got: %s", entryCSS, entries[0])
	}
}
