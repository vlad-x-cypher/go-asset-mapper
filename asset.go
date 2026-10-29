package asset

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var errInvalidHashLen = errors.New("hashLen must be between 0 and 64")

type Asset struct {
	PublicPath string
	Hash       string
	Path       string
}

func NewAsset(file *os.File, path, publicPath string, hashLen int) (*Asset, error) {
	defer file.Close()

	hashLen = max(0, hashLen)
	if hashLen > 64 {
		return nil, errInvalidHashLen
	}

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
		Path:       strings.TrimLeft(filepath.ToSlash(path), "/"),
		Hash:       hash,
		PublicPath: strings.TrimRight(publicPath, "/") + "/",
	}, nil
}

func (a *Asset) String() string {
	if a.Hash != "" {
		return a.PublicPath + a.Path + "?v=" + a.Hash
	}
	return a.PublicPath + a.Path
}
