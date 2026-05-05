package crypto

import "testing"

func TestEncryptDecrypt(t *testing.T) {

	s := New("secret")

	enc, err := s.Encrypt([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	dec, err := s.Decrypt(enc)
	if err != nil {
		t.Fatal(err)
	}

	if string(dec) != "hello" {
		t.Fatal("wrong decrypt")
	}
}
