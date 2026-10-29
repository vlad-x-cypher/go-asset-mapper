package asset

import "regexp"

var (
	cssRe   = regexp.MustCompile(`\.css$`)
	jsRe    = regexp.MustCompile(`\.js|\.mjs|\.cjs$`)
	imageRe = regexp.MustCompile(`(\.webp|\.jpg|\.jpeg|\.jpe|\.jfif|\.jif|\.png|\.gif|\.tiff|\.tif|\.svg|\.avif)$`)
)

func isCSS(path string) bool {
	return cssRe.MatchString(path)
}

func isJS(path string) bool {
	return jsRe.MatchString(path)
}

func isImage(path string) bool {
	return imageRe.MatchString(path)
}
