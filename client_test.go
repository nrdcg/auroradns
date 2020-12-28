package auroradns

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (*Client, *http.ServeMux) {
	t.Helper()

	apiHandler := http.NewServeMux()
	server := httptest.NewServer(apiHandler)

	client, err := NewClient(nil, WithBaseURL(server.URL))
	require.NoError(t, err)

	t.Cleanup(server.Close)

	return client, apiHandler
}

func handleAPI(mux *http.ServeMux, pattern, method string, next func(w http.ResponseWriter, r *http.Request)) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, fmt.Sprintf("invalid method %s", r.Method), http.StatusMethodNotAllowed)
			return
		}

		contentType := r.Header.Get(contentTypeHeader)
		if contentType != contentTypeJSON {
			http.Error(w, fmt.Sprintf("invalid Content-Type %s", contentType), http.StatusBadRequest)
			return
		}

		if next != nil {
			next(w, r)
		}
	})
}
