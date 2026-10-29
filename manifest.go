package asset

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

type ManifestType int

const (
	ViteManifestType ManifestType = iota
	WebpackManifestType
)

// ManifestConfig provides information about manifest filepath and how to parse it correctly.
// Currently supports Vite or Webpack types.
type ManifestConfig struct {
	// manifest filepath
	Path string
	// manifest generator type
	Type ManifestType
}

type viteManifestRecord struct {
	File           string   `json:"file"`
	Src            string   `json:"src"`
	Name           string   `json:"name"`
	IsEntry        bool     `json:"isEntry"`
	CSS            []string `json:"css"`
	Imports        []string `json:"imports"`
	IsDynamicEntry bool     `json:"isDynamicEntry"`
}

func parseManifest(path string, a *AssetMapper, t ManifestType) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	switch t {
	case ViteManifestType:
		return readViteManifest(file, a)
	case WebpackManifestType:
		return readWebpackManifest(file, a)
	}

	return errors.New("undefined manifest type")
}

func readViteManifest(reader io.Reader, a *AssetMapper) error {
	var err error
	decoder := json.NewDecoder(reader)

	for decoder.More() {
		var data map[string]viteManifestRecord

		err = decoder.Decode(&data)
		if err != nil {
			return err
		}

		for k, v := range data {
			asset, err := NewAsset(nil, v.File, a.PublicPath, 0)
			if err != nil {
				return err
			}

			a.AddAsset(asset, false)
			a.Assets[k] = asset
			if v.IsEntry {
				entry := a.CreateEntry(v.Name)
				entry.Add(asset.String())

				for _, css := range v.CSS {
					asset, err := NewAsset(nil, css, a.PublicPath, 0)
					if err != nil {
						return err
					}

					a.AddAsset(asset, false)
					entry.Add(asset.String())
				}
			}
		}
	}

	return nil
}

func readWebpackManifest(reader io.Reader, a *AssetMapper) (err error) {
	decoder := json.NewDecoder(reader)

	for decoder.More() {
		var data map[string]string

		err = decoder.Decode(&data)
		if err != nil {
			return err
		}

		for k, v := range data {
			asset, err := NewAsset(nil, v, a.PublicPath, 0)
			if err != nil {
				return err
			}
			a.AddAsset(asset, false)
			a.Assets[k] = asset
		}

	}
	return nil
}
