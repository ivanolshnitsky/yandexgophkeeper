package client

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status":"ok"}`))
		}),
	)
	defer server.Close()

	app := New(server.URL)

	resp, err := app.Ping()

	assert.NoError(t, err)
	assert.Contains(t, resp, "ok")
}

func TestSaveLoadToken(t *testing.T) {
	defer os.Remove(tokenFile)

	err := SaveToken("token123")
	assert.NoError(t, err)

	token, err := LoadToken()

	assert.NoError(t, err)
	assert.Equal(t, "token123", token)
}

func TestLoadToken_NotFound(t *testing.T) {
	_ = os.Remove(tokenFile)

	_, err := LoadToken()

	assert.Error(t, err)
}
