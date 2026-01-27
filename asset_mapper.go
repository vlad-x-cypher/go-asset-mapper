// Package asset provides asset mapper functionality to help using static assets (css, scripts, images)
// in Go templates
//
// Package inspired by Symfony AssetMapper, only much simpler, without importmap functionality and compiling assets.
package asset

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type AssetMapperEntry struct {
	CSS []string
	JS  []string
}

type AssetMapper struct {
	PublicPath string
	Assets     map[string]*Asset
	Entries    map[string]*AssetMapperEntry
	HashLen    int
	// Left trim subsrtring from final path
	Trim string
}

func NewAssetMapper() *AssetMapper {
	return &AssetMapper{
		Assets:     map[string]*Asset{},
		PublicPath: "/",
		HashLen:    10,
		Entries:    map[string]*AssetMapperEntry{},
		Trim:       "",
	}
}

// UseManifest loads all assets from provided manifest config.
// For more information look [ManifestConfig]
//
// Example:
//
//	assetMapper.UseManifest(&asset.ManifestConfig{
//		Path: "assets/manifest.json",
//		Type: asset.ViteManifestType,
//	})
func (a *AssetMapper) UseManifest(config ManifestConfig) error {
	switch config.Type {
	case ViteManifestType:
		return parseViteManifest(config.Path, a)
	case WebpackManifestType:
		return parseWebpackManifest(config.Path, a)
	}
	return errors.New("undefined manifest type")
}

// CreateEntry creates AssetsMapperEntry if not exists and returns pointer to that entry.
func (a *AssetMapper) CreateEntry(name string) *AssetMapperEntry {
	if e, ok := a.Entries[name]; ok {
		return e
	}

	a.Entries[name] = &AssetMapperEntry{
		CSS: []string{},
		JS:  []string{},
	}

	return a.Entries[name]
}

func (entry *AssetMapperEntry) Add(path string) {
	switch {
	case isCSS(path):
		entry.CSS = append(entry.CSS, path)
	case isJS(path):
		entry.JS = append(entry.JS, path)
	}
}

// AddAsset adds asset to list. If renew is set to true, existing asset will be
// replaced by provided one.
func (a *AssetMapper) AddAsset(asset *Asset, renew bool) {
	if !renew {
		if _, ok := a.Assets[asset.Path]; ok {
			return
		}
	}

	a.Assets[asset.Path] = asset
	a.Assets[asset.PublicPath] = asset
}

// ScanDir walks directory and maps all files to AssetMapper, storing its path and hash.
func (a *AssetMapper) ScanDir(dirName string) error {
	err := filepath.Walk(dirName, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		f, e := os.Open(path)

		if e != nil {
			return e
		}

		if a.Trim != "" {
			path = strings.TrimPrefix(path, a.Trim)
		}

		asset, assetErr := NewAsset(f, path, a.PublicPath, a.HashLen)
		if assetErr != nil {
			return assetErr
		}

		a.AddAsset(asset, false)

		return nil
	})

	return err
}

func extractAssetPathFromMap(m map[string]*Asset, search string) string {
	search = strings.TrimLeft(search, "/")
	if asset, ok := m[search]; ok {
		return asset.String()
	}
	return search
}

// Get returns asset url including version. If asset not found returns path param as is.
func (a *AssetMapper) Get(path string) string {
	return extractAssetPathFromMap(a.Assets, path)
}

func attributeMapToString(m map[string]string) string {
	s := []string{}

	for k, v := range m {
		if v == "" {
			s = append(s, html.EscapeString(k))
		} else {
			s = append(s, fmt.Sprintf(`%s="%s"`, html.EscapeString(k), html.EscapeString(v)))
		}
	}

	return strings.Join(s, " ")
}

// tagAttributes Extracts attributes to key: val map from flat slice
// attrs must be an even number, every first reperesenting key and every second - value
// Example: ["attr1", "val1", "attr2", "val2"]
//
// Result: map[string]string{"attr1": "val1", "attr2": "val2"}
func tagAttributes(attrs []string) (map[string]string, error) {
	if len(attrs)%2 != 0 {
		return nil, errors.New("attrs must be an even number of strings")
	}
	attrMap := map[string]string{}

	for i := 0; i < len(attrs); i += 2 {
		attrMap[attrs[i]] = attrs[i+1]
	}

	return attrMap, nil
}

func scriptTag(attrs string) template.HTML {
	return template.HTML(fmt.Sprintf("<script %s></script>", attrs))
}

// ScriptTag returns HTML script tag
// attrs param can be used to pass additional attributes to the tag. Unfortunately Go html.Template
// does not allow create new maps inside html templates, attrs must be an even number of strings
// reperesenting key value pairs.
//
// Example usage in template:
//
//	{{ scriptTag "main.js" }}
//
//	<!-- Passing additional attributes to script tag -->
//	{{ scriptTag "other.js" "type" "module" "id" "other-script" }}
//
//	<!-- Example set defer or async attributes -->
//	{{ scriptTag "defered.js" "defer" "" }}
//	{{ scriptTag "some-async.js" "async" "" }}
//
// Result:
//
//	<script src="main.js"></script>
//
//	<!-- Passing additional attributes to script tag -->
//	<script src="other.js" type="module" id="other-script"></script>
//
//	<!-- Example set defer or async attributes -->
//	<script defer src="defered.js"></script>
//	<script async src="some-async.js"></script>
func (a *AssetMapper) ScriptTag(path string, attrs ...string) (template.HTML, error) {
	link := a.Get(path)

	attrMap, err := tagAttributes(attrs)
	if err != nil {
		return "", err
	}

	attrMap["src"] = link

	return scriptTag(attributeMapToString(attrMap)), nil
}

func linkTag(attrs string) template.HTML {
	return template.HTML(fmt.Sprintf("<link %s/>", attrs))
}

// LinkTag returns HTML link tag
// attrs param can be used to pass additional attributes to the tag. Unfortunately Go html.Template
// does not allow create new maps inside html templates, attrs must be an even number of strings
// reperesenting key value pairs.
//
// Example usage in template:
//
//	{{ linkTag "style.css" }}
//
//	<!-- Passing additional attributes to link tag -->
//	{{ linkTag "homepage.css" "id" "homepage-css" "media" "screen" }}
//
// Result:
//
//	<link href="style.css" rel="stylesheet"/>
//	<!-- Passing additional attributes to link tag -->
//	<link href="homepage.css" rel="stylesheet" id="homepage-css" media="screen"/>
func (a *AssetMapper) LinkTag(path string, attrs ...string) (template.HTML, error) {
	link := a.Get(path)

	attrs = append([]string{"rel", "stylesheet"}, attrs...)
	attrMap, err := tagAttributes(attrs)
	if err != nil {
		return "", err
	}

	attrMap["href"] = link

	return linkTag(attributeMapToString(attrMap)), nil
}

// CSSEntry returns slice of css urls from entrypoint
func (a *AssetMapper) CSSEntry(name string) []string {
	if s, ok := a.Entries[name]; ok {
		return s.CSS
	}
	return nil
}

// JSEntry returns slice of js urls from entrypoint
func (a *AssetMapper) JSEntry(name string) []string {
	if s, ok := a.Entries[name]; ok {
		return s.JS
	}
	return nil
}

// CSSLinkTagsFromEntry return slice of html links from entry.
//
// For more information look [AssetMapper.LinkTag] method
func (a *AssetMapper) CSSLinkTagsFromEntry(name string, attrs ...string) ([]template.HTML, error) {
	attrs = append([]string{"rel", "stylesheet"}, attrs...)
	attrMap, err := tagAttributes(attrs)
	if err != nil {
		return nil, err
	}

	result := []template.HTML{}
	for _, css := range a.CSSEntry(name) {
		attrMap["href"] = css
		result = append(result, linkTag(attributeMapToString(attrMap)))
	}

	return result, nil
}

// JSScriptTagsFromEntry return slice of html scripts from entry.
//
// For more information look [AssetMapper.ScriptTag] method
func (a *AssetMapper) JSScriptTagsFromEntry(name string, attrs ...string) ([]template.HTML, error) {
	attrMap, err := tagAttributes(attrs)
	if err != nil {
		return nil, err
	}

	result := []template.HTML{}
	for _, js := range a.JSEntry(name) {
		attrMap["src"] = js
		result = append(result, scriptTag(attributeMapToString(attrMap)))
	}

	return result, nil
}
