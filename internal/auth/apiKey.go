package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	keyString := headers.Get("Authorization")
	if keyString == "" {
		return "", fmt.Errorf("authorization not found")
	}
	s := strings.Fields(keyString)
	return s[1], nil
}
