package asset

import "testing"

func TestIsCSS(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"simple", "style.css", true},
		{"nested", "assets/css/style.css", true},
		{"extension only", ".css", true},
		{"css after js", "style.js.css", true},
		{"scss", "style.scss", false},
		{"source map", "style.css.map", false},
		{"js", "style.js", false},
		{"no extension", "css", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCSS(tt.path); got != tt.want {
				t.Errorf("isCSS(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsJS(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"js", "app.js", true},
		{"mjs", "app.mjs", true},
		{"cjs", "app.cjs", true},
		{"nested", "assets/js/app.mjs", true},
		{"extension only", ".js", true},
		{"json", "app.json", false},
		{"source map", "app.js.map", false},
		{"mjs backup", "app.mjs.bak", false},
		{"css", "app.css", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isJS(tt.path); got != tt.want {
				t.Errorf("isJS(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
