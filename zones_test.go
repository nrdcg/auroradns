package auroradns

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateZoneWithContext(t *testing.T) {
	client, mux := setupTest(t)

	handleAPI(mux, "/zones", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if string(reqBody) != `{"name":"example.com"}` {
			http.Error(w, fmt.Sprintf("invalid request body: %s", string(reqBody)), http.StatusInternalServerError)
			return
		}

		_, err = fmt.Fprintf(w, `
				{
					"id":   "identifier-zone-1",
					"name": "example.com"
				}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	zone, resp, err := client.CreateZoneWithContext(t.Context(), "example.com")
	require.NoError(t, err)

	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expected := &Zone{ID: "identifier-zone-1", Name: "example.com"}
	assert.Equal(t, expected, zone)
}

func TestClient_CreateZoneWithContext_error(t *testing.T) {
	client, mux := setupTest(t)

	handleAPI(mux, "/zones", http.MethodPost, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)

		_, err := fmt.Fprintf(w, `{
  			"error": "AuthenticationRequiredError",
  			"errormsg": "Failed to parse Authorization header"
			}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	zone, resp, err := client.CreateZoneWithContext(t.Context(), "example.com")
	require.Error(t, err)

	require.NotNil(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	assert.Nil(t, zone)
}

func TestClient_DeleteZoneWithContext(t *testing.T) {
	client, mux := setupTest(t)

	handleAPI(mux, "/zones/identifier-zone-1", http.MethodDelete, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	result, resp, err := client.DeleteZoneWithContext(t.Context(), "identifier-zone-1")
	require.NoError(t, err)

	require.NotNil(t, resp)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	assert.True(t, result)
}

func TestClient_DeleteZoneWithContext_error(t *testing.T) {
	client, mux := setupTest(t)

	handleAPI(mux, "/zones/identifier-zone-1", http.MethodDelete, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)

		_, err := fmt.Fprintf(w, `{
  			"error": "AuthenticationRequiredError",
  			"errormsg": "Failed to parse Authorization header"
			}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	result, resp, err := client.DeleteZoneWithContext(t.Context(), "identifier-zone-1")
	require.Error(t, err)

	require.NotNil(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	assert.False(t, result)
}

func TestClient_ListZonesWithContext(t *testing.T) {
	client, mux := setupTest(t)

	handleAPI(mux, "/zones", http.MethodGet, func(w http.ResponseWriter, _ *http.Request) {
		_, err := fmt.Fprintf(w, `[
				{
					"id":   "identifier-zone-1",
					"name": "example.com"
				}
			]`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	zones, resp, err := client.ListZonesWithContext(t.Context())
	require.NoError(t, err)

	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expected := []Zone{{ID: "identifier-zone-1", Name: "example.com"}}
	assert.Equal(t, expected, zones)
}

func TestClient_ListZonesWithContext_error(t *testing.T) {
	client, mux := setupTest(t)

	handleAPI(mux, "/zones", http.MethodGet, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)

		_, err := fmt.Fprintf(w, `{
  			"error": "AuthenticationRequiredError",
  			"errormsg": "Failed to parse Authorization header"
			}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	zones, resp, err := client.ListZonesWithContext(t.Context())
	require.EqualError(t, err, "AuthenticationRequiredError - Failed to parse Authorization header")

	require.NotNil(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	assert.Nil(t, zones)
}
