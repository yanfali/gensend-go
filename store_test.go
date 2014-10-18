package main

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
)

func getStoreUrl() string {
	return test_server.URL + "/store"
}

func getClientPostResponse(t *testing.T, content []byte) *http.Response {
	req, err := http.NewRequest("POST", getStoreUrl(), bytes.NewBuffer(content))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	return resp
}

func TestStoreOK(t *testing.T) {
	var jsonStr = []byte(`{"password":"abc123", "maxReads": 1, "maxMinutes": 1}`)

	resp := getClientPostResponse(t, jsonStr)
	defer resp.Body.Close()
	assertStatusCode(t, resp, http.StatusOK)

	//decoder := json.NewDecoder(resp.Body)
	var jsonResp JSONStoreResponse
	assertValidJSON(t, resp, &jsonResp)

	if !strings.HasSuffix(jsonResp.Url, jsonResp.Token) {
		t.Logf("expected Url to end with %q, got %q", jsonResp.Token, jsonResp.Url)
		t.Fail()
	}
}

func TestStoreBadPassword(t *testing.T) {
	var jsonStr = []byte(`{"password":"", "maxReads": 1, "maxMinutes": 1}`)
	resp := getClientPostResponse(t, jsonStr)
	defer resp.Body.Close()
	assertStatusCode(t, resp, http.StatusBadRequest)

	var jsonResp JSONErrorResponse
	assertValidJSON(t, resp, &jsonResp)

	if jsonResp.StatusCode != http.StatusBadRequest || !strings.Contains(jsonResp.ErrorMessage, "Password") {
		t.Logf("Expected Password Error, got %q", jsonResp.ErrorMessage)
		t.Fail()
	}
}
