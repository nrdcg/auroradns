package auroradns

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateZone(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	handleAPI(mux, "/zones", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
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

	zone, resp, err := client.CreateZone("example.com")
	require.NoError(t, err)

	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expected := &Zone{ID: "identifier-zone-1", Name: "example.com"}
	assert.Equal(t, expected, zone)
}

func TestClient_CreateZone_error(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	handleAPI(mux, "/zones", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
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

	zone, resp, err := client.CreateZone("example.com")
	require.Error(t, err)

	require.NotNil(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	assert.Nil(t, zone)
}

func TestClient_DeleteZone(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	handleAPI(mux, "/zones/identifier-zone-1", http.MethodDelete, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	result, resp, err := client.DeleteZone("identifier-zone-1")
	require.NoError(t, err)

	require.NotNil(t, resp)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	assert.True(t, result)
}

func TestClient_DeleteZone_error(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	handleAPI(mux, "/zones/identifier-zone-1", http.MethodDelete, func(w http.ResponseWriter, r *http.Request) {
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

	result, resp, err := client.DeleteZone("identifier-zone-1")
	require.Error(t, err)

	require.NotNil(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	assert.False(t, result)
}

func TestClient_ListZones(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	handleAPI(mux, "/zones", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
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

	zones, resp, err := client.ListZones()
	require.NoError(t, err)

	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expected := []Zone{{ID: "identifier-zone-1", Name: "example.com"}}
	assert.Equal(t, expected, zones)
}

func TestClient_ListZones_error(t *testing.T) {
	client, mux, tearDown := setupTest()
	defer tearDown()

	handleAPI(mux, "/zones", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
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

	zones, resp, err := client.ListZones()
	require.EqualError(t, err, "AuthenticationRequiredError - Failed to parse Authorization header")

	require.NotNil(t, resp)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	assert.Nil(t, zones)
}
