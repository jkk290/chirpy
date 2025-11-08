package auth

import (
	"net/http"
	"testing"
)

func TestBearerToken(t *testing.T) {
	tokenString := "someTokenString"
	header := http.Header{}
	header.Add("Authorization", "Bearer someTokenString")
	token, err := GetBearerToken(header)
	if err != nil || token == "" {
		t.Errorf("error getting token: %s", err)
		t.Fail()
	}
	if token != tokenString {
		t.Fail()
	}
}

func TestNoBearerToken(t *testing.T) {
	header := http.Header{}
	token, err := GetBearerToken(header)
	if err == nil || token != "" {
		t.Fail()
	}
}
