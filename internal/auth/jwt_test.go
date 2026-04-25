package auth

import "testing"

func TestJWT(t *testing.T) {
	j := NewJWTService("secret")

	token, err := j.Generate("user")
	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("empty token")
	}
}
