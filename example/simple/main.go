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

	err := assetMapper.ScanDir("assets")
	if err != nil {
		log.Fatalf("assets scandir err: %v", err)
	}

	t.Funcs(template.FuncMap{
		"scriptTag": assetMapper.ScriptTag,
		"linkTag":   assetMapper.LinkTag,
		"asset":     assetMapper.Get,
	})

	templates, err := t.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "home", nil)
		if err != nil {
			log.Fatal(err)
		}
	})
	http.Handle("GET /assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
