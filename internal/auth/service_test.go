package auth

import "testing"

func TestRegisterAndLogin(t *testing.T) {
	s := NewService()

	err := s.Register("user", "pass")
	if err != nil {
		t.Fatal(err)
	}

	err = s.Login("user", "pass")
	if err != nil {
		t.Fatal("login should succeed")
	}
}

func TestRegisterDuplicate(t *testing.T) {
	s := NewService()

	_ = s.Register("user", "pass")
	err := s.Register("user", "pass")

	if err == nil {
		t.Fatal("expected error for duplicate user")
	}
}
