package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/nit-app/nit-backend/models/requests"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

func makeClient(t *testing.T) *resty.Client {
	testPhoneNumber := os.Getenv("MOCK_PHONE_NUMBER")
	testCode := os.Getenv("OTP_MOCK_ACCEPT_CODE")
	endpoint := os.Getenv("ENDPOINT")

	require.NotEmpty(t, endpoint, "endpoint should not be empty")
	require.NotEmpty(t, testPhoneNumber, "phone number should not be empty")
	require.NotEmpty(t, testCode, "code should not be empty")

	b, err := json.Marshal(requests.PhoneNumberRequest{PhoneNumber: testPhoneNumber})
	require.NoError(t, err)

	resp, err := http.Post(endpoint+"/auth/sendCode", "application/json", bytes.NewReader(b))
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	b, err = json.Marshal(requests.OtpCheckRequest{Code: testCode})
	require.NoError(t, err)

	req, err := http.NewRequest("POST", endpoint+"/auth/confirm", bytes.NewReader(b))
	require.NoError(t, err)

	for _, c := range resp.Cookies() {
		req.AddCookie(c)
	}

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	client := resty.New()
	client.SetCookies(req.Cookies())
	client.SetBaseURL(endpoint)

	return client
}
