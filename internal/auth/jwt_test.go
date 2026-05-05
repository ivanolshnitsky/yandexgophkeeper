package auth

import (
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	j := NewJWTService("secret", time.Hour)

	token, err := j.Generate("user")
	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("empty token")
	}
}
