package internal

import (
	"io"
	"net/http"
)

func GetPublicIP() (string, error) {
	resp, err := http.Get("http://ifconfig.me")
	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
