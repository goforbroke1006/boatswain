package internal

import (
	"io/ioutil"
	"net/http"
)

func GetPublicIP() (string, error) {
	resp, err := http.Get("http://ifconfig.me")
	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
