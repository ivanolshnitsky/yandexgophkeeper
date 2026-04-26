package storage

import "testing"

func TestMemoryStorage(t *testing.T) {
	s := NewMemory()

	s.Save("user", Credential{Login: "l", Password: "p"})

	data := s.List("user")

	if len(data) != 1 {
		t.Fatal("expected 1 item")
	}
}
