package auroradns

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTokenTransport_success(t *testing.T) {
	apiKey := "☺"
	secretKey := "🔑"

	transport, err := NewTokenTransport(apiKey, secretKey)
	require.NoError(t, err)
	assert.NotNil(t, transport)
}

func TestNewTokenTransport_missing_credentials(t *testing.T) {
	apiKey := ""
	secretKey := ""

	transport, err := NewTokenTransport(apiKey, secretKey)
	require.Error(t, err)
	assert.Nil(t, transport)
}

func TestTokenTransport_RoundTrip(t *testing.T) {
	apiKey := "☺"
	secretKey := "🔑"

	transport, err := NewTokenTransport(apiKey, secretKey)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	resp, err := transport.RoundTrip(req)
	require.NoError(t, err)

	assert.Regexp(t, `\d{8}T\d{6}Z`, resp.Request.Header.Get("X-Auroradns-Date"))
	assert.Regexp(t, `AuroraDNSv1 \w{64}`, resp.Request.Header.Get("Authorization"))
}
