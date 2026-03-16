package utils

import (
	"crypto/tls"
	"io"
	"net/http"
)

func GetHttpResponse(url string) (body []byte, err error) {
	var (
		req *http.Request
		res *http.Response
	)
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	res, err = httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	return
}
