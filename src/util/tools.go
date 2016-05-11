package util

import (
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
)

import (
	"golang.org/x/net/html/charset"
)

func IsFileExist(name string) bool {
	if name == "" {
		return false
	}
	fi, err := os.Stat(name)
	if err != nil {
		return os.IsExist(err)
	} else {
		return !fi.IsDir()
	}
}

func IsDirExist(path string) bool {
	if path == "" {
		return false
	}
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

func Mkdir(dirpath string) (bool, error) {
	if dirpath == "" {
		return false, errors.New("dirpath is empty")
	}
	dirpath = path.Clean(dirpath)
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		return false, err
	}
	return true, nil
}

func RelativeUrl2Abs(baseurl *url.URL, refurl string) (string, error) {
	if refurl == "" {
		return "", errors.New("refurl is empty")
	}
	if baseurl == nil {
		return "", errors.New("baseurl is nil")
	}

	rel, err := url.Parse(refurl)
	if err != nil {
		return "", err
	}
	if rel.IsAbs() {
		return refurl, nil
	}
	return baseurl.ResolveReference(rel).String(), nil
}

func GetUrlFromJs(js string) string {
	if strings.HasPrefix(js, "javascript:location.href") {
		url := strings.Replace(js, "javascript:location.href=", "", 1)
		url = strings.Trim(url, "\"'")
		return url
	}
	return js
}

func Charset2UTF8(contentType string, reader io.Reader) (string, error) {
	destReader, err := charset.NewReader(reader, contentType)
	if err != nil {
		return "", err
	}

	var targetBytes []byte
	if targetBytes, err = ioutil.ReadAll(destReader); err != nil {
		return "", err
	}
	return string(targetBytes), nil
}
