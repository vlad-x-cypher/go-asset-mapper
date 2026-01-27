package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/vlad-x-cypher/go-asset-mapper"
)

func main() {
	t := template.New("")

	assetMapper := asset.NewAssetMapper()

	assetMapper.PublicPath = "/static/"

	err := assetMapper.UseManifest(asset.ManifestConfig{
		Path: "public/.vite/manifest.json",
		Type: asset.ViteManifestType,
	})
	if err != nil {
		log.Fatalf("assets vite manifest parse err: %v", err)
	}

	err = assetMapper.UseManifest(asset.ManifestConfig{
		Path: "public/bundle/manifest.json",
		Type: asset.WebpackManifestType,
	})
	if err != nil {
		log.Fatalf("assets webpack manifest parse err: %v", err)
	}

	t.Funcs(template.FuncMap{
		"scriptTag":      assetMapper.ScriptTag,
		"linkTag":        assetMapper.LinkTag,
		"asset":          assetMapper.Get,
		"entryCss":       assetMapper.CSSEntry,
		"entryCssLinks":  assetMapper.CSSLinkTagsFromEntry,
		"entryJsScripts": assetMapper.JSScriptTagsFromEntry,
	})

	templates, err := t.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "index", nil)
		if err != nil {
			log.Fatal(err)
		}
	})
	http.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./public"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
