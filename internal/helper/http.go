package helper

import (
	"fmt"
	"net/http"
)

func Get(url string, header http.Header) (*http.Response, error) {
	client := http.DefaultClient

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error %s", err.Error())
	}
	req.Header = header

	return client.Do(req)
}
