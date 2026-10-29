package asset

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/url"
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
	if file != nil {
		defer func() {
			_ = file.Close()
		}()
	}

	hashLen = max(0, hashLen)
	if hashLen > 64 {
		return nil, errInvalidHashLen
	}

	hash := ""
	if hashLen > 0 && !isAbsoluteURL(path) && file != nil {
		hasher := sha256.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return nil, err
		}

		hash = hex.EncodeToString(hasher.Sum(nil))
		hash = hash[0:hashLen]
	}

	if isAbsoluteURL(path) {
		publicPath = ""
	} else {
		publicPath = strings.TrimRight(publicPath, "/") + "/"
		path = strings.TrimLeft(filepath.ToSlash(path), "/")
	}

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

func isAbsoluteURL(s string) bool {
	if strings.HasPrefix(s, "//") {
		return true
	}
	u, err := url.Parse(s)
	return err == nil && u.IsAbs()
}
