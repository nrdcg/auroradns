package auroradns

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTokenTransport_success(t *testing.T) {
	userID := "☺"
	key := "🔑"

	transport, err := NewTokenTransport(userID, key)
	require.NoError(t, err)
	assert.NotNil(t, transport)
}

func TestNewTokenTransport_missing_credentials(t *testing.T) {
	userID := ""
	key := ""

	transport, err := NewTokenTransport(userID, key)
	require.Error(t, err)
	assert.Nil(t, transport)
}

func TestTokenTransport_RoundTrip(t *testing.T) {
	userID := "☺"
	key := "🔑"

	transport, err := NewTokenTransport(userID, key)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	_, err = transport.RoundTrip(req)
	require.NoError(t, err)

	assert.Regexp(t, `\d{8}T\d{6}Z`, req.Header.Get("X-Auroradns-Date"))
	assert.Regexp(t, `AuroraDNSv1 \w{64}`, req.Header.Get("Authorization"))
}
