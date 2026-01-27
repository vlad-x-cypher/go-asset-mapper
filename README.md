# Go Asset Mapper

go-asset-mapper is a Go package to serve versioned static files in Go HTML Templates. Main feature: map and version all static files in specific directory.

## Instalation

Install via Go modules:
```bash
go get github.com/vlad-x-cypher/go-asset-mapper
```

### Example 

```go
 
package main

import (
	"html/template"
	"log"
	"os"

	"github.com/vlad-x-cypher/go-asset-mapper"
)

func main() {
	t := template.New("")

	assetMapper := asset.NewAssetMapper()

	// Scan specific directory to map all files
	err := assetMapper.ScanDir("assets")
	if err != nil {
		log.Fatalf("assets scandir err: %v", err)
	}

	// Add functions to use inside templates
	t.Funcs(template.FuncMap{
		"scriptTag": assetMapper.ScriptTag,
		"linkTag":   assetMapper.LinkTag,
		"asset":     assetMapper.Get,
	})

    const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		{{ linkTag "assets/style.css" }}
		{{ scriptTag "assets/main.js" }}
	</head>
	<body>
		<img src="{{ asset "assets/example.png" }}" />
	</body>
</html>`

    templates, err := t.Parse(tpl)
    if err != nil {
        log.Fatal(err)
    }

    err = templates.Execute(os.Stdout, nil)
    if err != nil {
        log.Fatal(err)
    }
}
```
Output:
```html
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<link href="assets/style.css" rel="stylesheet"/>
		<script src="assets/main.js"></script>
	</head>
	<body>
		<img src="assets/example.png" />
	</body>
</html>
```

## Map assets from manifest.json

Javascript Build tools (Vite, Webpack) allows to generate `manifest.json` with list of all processed assets. 

`AssetsMapper.UseManifest` method designed to map all assets from manifest.

### Example

example manifest.json:
```json
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
}
```

```go
 
package main

import (
	"html/template"
	"log"
	"os"

	"github.com/vlad-x-cypher/go-asset-mapper"
)

func main() {
	t := template.New("")

	assetMapper := asset.NewAssetMapper()

    	// Parse manifest.json 
	err := assetMapper.UseManifest(&asset.ManifestConfig{
        	Path: "assets/.vite/manifest.json",
        	// Set correct type
		Type: asset.ViteManifestType,
    	})
	if err != nil {
		log.Fatalf("assets vite manifest err: %v", err)
	}

	// Add functions to use inside templates
	t.Funcs(template.FuncMap{
		"scriptTag": assetMapper.ScriptTag,
		"linkTag":   assetMapper.LinkTag,
		"asset":     assetMapper.Get,
		"entryCssLinks":  assetMapper.CSSLinkTagsFromEntry,
		"entryJsScripts": assetMapper.JSScriptTagsFromEntry,
	})

    const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
        {{ range (entryCssLinks "app") }}
            {{ . }}
        {{ end }}
        {{ range (entryJsScripts "app" "defer" "") }}
            {{ . }}
        {{ end }}
	</head>
	<body>
		<img src="{{ asset "src/placeholder.jpeg" }}" />
	</body>
</html>`

    templates, err := t.Parse(tpl)
    if err != nil {
        log.Fatal(err)
    }

    err = templates.Execute(os.Stdout, nil)
    if err != nil {
        log.Fatal(err)
    }
}
```
Output:
```html
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<link href="assets/app-o2N34dPp.css" rel="stylesheet"/>
		<script src="assets/app-CKgRTByK.js" defer></script>
	</head>
	<body>
		<img src="assets/placeholder-DXdl7YkJ.jpeg" />
	</body>
</html>
```

For more complete examples see [example dir](./example)
