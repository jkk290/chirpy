package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidJWTS(t *testing.T) {
	userId := uuid.New()
	tokenSecret := "ThisIsASecret123!"

	tokenString, err := MakeJWT(userId, tokenSecret, 1*time.Minute)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	userUUID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if userId != userUUID {
		t.Errorf("no match, %v, %v", userId, userUUID)
		t.Fail()
	}
}

func TestInvalidJWTS(t *testing.T) {
	userId := uuid.New()
	tokenSecret := "anotherSecret"
	wrongTokenSecret := "thisShouldbeSUPERwrong"

	tokenString, err := MakeJWT(userId, tokenSecret, 1*time.Minute)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	userUUID, err := ValidateJWT(tokenString, wrongTokenSecret)
	if err == nil || userUUID != uuid.Nil {
		t.Fail()
	}
}

func TestExpiredJWTS(t *testing.T) {
	userId := uuid.New()
	tokenSecret := "anotherSecret"

	tokenString, _ := MakeJWT(userId, tokenSecret, 1*time.Second)
	time.Sleep(1200 * time.Millisecond)
	_, err := ValidateJWT(tokenString, tokenSecret)
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fail()
	}
}
