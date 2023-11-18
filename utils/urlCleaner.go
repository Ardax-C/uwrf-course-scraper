package utils

import (
	"net/url"
	"strings"
)

func CleanURL(rawLink string) (string, error) {
	parsedURL, err := url.Parse(rawLink)
	if err != nil {
		return "", err
	}

	// Clean each query parameter
	q := parsedURL.Query()
	for key := range q {
		q.Set(key, strings.TrimSpace(q.Get(key)))
	}
	parsedURL.RawQuery = q.Encode()

	return parsedURL.String(), nil
}
