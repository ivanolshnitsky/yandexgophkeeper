package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	app := New(server.URL)

	resp, err := app.Ping()

	assert.NoError(t, err)
	assert.Contains(t, resp, "ok")
}
