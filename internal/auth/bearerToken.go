package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("authorization not found")
	}
	splitted := strings.Split(tokenString, "Bearer ")
	return splitted[0], nil
}
