package asset

import (
	"testing"
)

func TestAssetMapperGet(t *testing.T) {
	a := NewAssetMapper()
	a.Assets["test.css"] = &Asset{
		Path:       "test.css",
		Hash:       "123",
		PublicPath: "/",
	}

	table := []struct {
		in       string
		expected string
	}{
		{"test.css", "/test.css?v=123"},
		{"raw.doc", "raw.doc"},
		{"/raw.doc", "/raw.doc"},
	}

	for _, tt := range table {
		result := a.Get(tt.in)
		if tt.expected != result {
			t.Errorf("String should be equal. Expected: %s\nGot:%s\n", tt.expected, result)
		}
	}
}

func TestAssetMapperGetAbsoluteURL(t *testing.T) {
	table := []struct {
		in       string
		expected string
	}{
		{"//example.com/path.js", "//example.com/path.js"},
		{"http://example.com/path.js", "http://example.com/path.js"},
		{"https://example.com/path.js", "https://example.com/path.js"},
	}
	a := NewAssetMapper()
	a.PublicPath = "/static"
	for _, v := range table {
		asset, err := NewAsset(nil, v.in, a.PublicPath, 32)
		if err != nil {
			t.Error(err)
		}
		a.AddAsset(asset, false)
	}
	for _, tt := range table {
		result := a.Get(tt.in)
		if tt.expected != result {
			t.Errorf("String should be equal. Expected: %s\nGot:%s\n", tt.expected, result)
		}
	}
}

func TestAttributeToString(t *testing.T) {
	attrTest := []struct {
		in       map[string]string
		expected string
	}{
		{
			map[string]string{
				"data-test": "value",
			},
			"data-test=\"value\"",
		},
		{
			map[string]string{
				"shouldEscape<>": ">",
			},
			"shouldEscape&lt;&gt;=\"&gt;\"",
		},
		{
			map[string]string{
				"empty": "",
			},
			"empty",
		},
	}

	for _, tt := range attrTest {
		s := attributeMapToString(tt.in)
		if s != tt.expected {
			t.Errorf("String should be equal. Expected: \"%s\"\nGot: \"%s\"\n", tt.expected, s)
		}
	}
}
