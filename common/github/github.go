package github

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func GetFileByToken(token, user, repo, branch, filePath string) (io.Reader, error) {
	filePath = strings.TrimPrefix(filePath, "/")
	fileUrl := fmt.Sprintf("https://%s@raw.githubusercontent.com/%s/%s/%s/%s", token, user, repo, branch, filePath)
	if token == "" {
		fileUrl = fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", user, repo, branch, filePath)
	}
	u, err := url.Parse(fileUrl)
	if err != nil {
		return nil, errors.Wrap(err, "parse url")
	}
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, errors.Wrap(err, "get file")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("status code: %d", resp.StatusCode)
	}
	var result = bytes.NewBuffer(nil)
	_, err = result.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body")
	}
	return result, nil
}
