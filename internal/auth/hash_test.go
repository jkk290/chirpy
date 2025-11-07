package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "HelloPASS435$"
	hashedPassword, _ := HashPassword(password)
	match, _ := CheckPasswordHash(password, hashedPassword)

	if !match {
		t.Fail()
	}
}

func TestWrongHashPassword(t *testing.T) {
	password := "FHgnda325#!ckdhH*."
	hashedPassword, _ := HashPassword(password)
	match, _ := CheckPasswordHash("123456!", hashedPassword)

	if match {
		t.Fail()
	}
}
