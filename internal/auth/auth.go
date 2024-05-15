package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Get API key from the headers of the request
// Example:
// Authorization: ApiKey {api key}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication found")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("invalid auth header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("invalid first part of Auth header")
	}
	return vals[1], nil
}