package asset

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

type Asset struct {
	PublicPath string
	Hash       string
	Path       string
}

func NewAsset(file *os.File, path, publicPath string, hashLen int) (*Asset, error) {
	defer file.Close()

	hash := ""
	if hashLen > 0 {
		hasher := sha256.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return nil, err
		}

		hash = hex.EncodeToString(hasher.Sum(nil))
	}

	hash = hash[0:hashLen]

	return &Asset{
		Path:       path,
		Hash:       hash,
		PublicPath: publicPath,
	}, nil
}

func (a *Asset) String() string {
	if a.Hash != "" {
		return a.PublicPath + a.Path + "?v=" + a.Hash
	}
	return a.PublicPath + a.Path
}
