package main

import (
	"net/http"
	"strings"
	"testing"
)

func getStoreUrl() string {
	return test_server.URL + getBaseAPIUrl() + "/store"
}

func JsonErrorsWrapper(t *testing.T, jsonStr []byte, expectedStatusCode int, fn func(jsonResp JSONErrorResponse)) {
	// Test Expecting An Error
	resp := getClientResponse(t, POST, jsonStr)
	defer resp.Body.Close()
	assertStatusCodeEquals(t, resp, expectedStatusCode)

	var jsonResp JSONErrorResponse
	assertValidJSON(t, resp, &jsonResp)
	fn(jsonResp)
}

func TestStoreOK(t *testing.T) {
	var jsonStr = []byte(`{"password":"abc123", "maxReads": 1, "maxMinutes": 1}`)

	resp := getClientResponse(t, POST, jsonStr)
	defer resp.Body.Close()
	assertStatusCodeEquals(t, resp, http.StatusOK)

	var jsonResp JSONStoreResponse
	assertValidJSON(t, resp, &jsonResp)

	if !strings.HasSuffix(jsonResp.Url, jsonResp.Token) {
		t.Logf("expected Url to end with %q, got %q", jsonResp.Token, jsonResp.Url)
		t.Fail()
	}
}

func TestUniqueTokenReturnedPerStore(t *testing.T) {
	var jsonStr = []byte(`{"password":"abc123", "maxReads": 1, "maxMinutes": 1}`)

	resp := getClientResponse(t, POST, jsonStr)
	assertStatusCodeEquals(t, resp, http.StatusOK)

	var jsonResp1 JSONStoreResponse
	assertValidJSON(t, resp, &jsonResp1)
	resp.Body.Close()

	resp = getClientResponse(t, POST, jsonStr)
	defer resp.Body.Close()
	assertStatusCodeEquals(t, resp, http.StatusOK)
	var jsonResp2 JSONStoreResponse
	assertValidJSON(t, resp, &jsonResp2)

	if jsonResp2.Token == jsonResp1.Token {
		t.Logf("Tokens unexpectedly matched for password %q, first %q, second %q", "abc123", jsonResp1.Token, jsonResp2.Token)
		t.Fail()
	}
}

func TestStoreInvalidPassword(t *testing.T) {
	JsonErrorsWrapper(t, []byte(`{"password":"", "maxReads": 1, "maxMinutes": 1}`), http.StatusBadRequest,
		func(jsonResp JSONErrorResponse) {
			if jsonResp.StatusCode != http.StatusBadRequest || !strings.Contains(jsonResp.ErrorMessage, "Password") {
				t.Logf("Expected Password Error, got %q", jsonResp.ErrorMessage)
				t.Fail()
			}
		})
}

func TestStoreInvalidMaxReads(t *testing.T) {
	JsonErrorsWrapper(t, []byte(`{"password":"321cba", "maxReads": 0, "maxMinutes": 1}`), http.StatusBadRequest,
		func(jsonResp JSONErrorResponse) {
			if jsonResp.StatusCode != http.StatusBadRequest || !strings.Contains(jsonResp.ErrorMessage, "maxReads") {
				t.Logf("Expected maxReads Error, got %q", jsonResp.ErrorMessage)
				t.Fail()
			}
		})
}

func TestStoreInvalidMaxReadsNegative(t *testing.T) {
	JsonErrorsWrapper(t, []byte(`{"password":"321cba", "maxReads": -1, "maxMinutes": 1}`), http.StatusBadRequest,
		func(jsonResp JSONErrorResponse) {
			if jsonResp.StatusCode != http.StatusBadRequest || !strings.Contains(jsonResp.ErrorMessage, "maxReads") {
				t.Logf("Expected maxReads Error, got %q", jsonResp.ErrorMessage)
				t.Fail()
			}
		})
}

func TestStoreInvalidMaxMinutes(t *testing.T) {
	JsonErrorsWrapper(t, []byte(`{"password":"321cba", "maxReads": 1, "maxMinutes": 0}`), http.StatusBadRequest,
		func(jsonResp JSONErrorResponse) {
			if jsonResp.StatusCode != http.StatusBadRequest || !strings.Contains(jsonResp.ErrorMessage, "maxMinutes") {
				t.Logf("Expected maxMinutes Error, got %q", jsonResp.ErrorMessage)
				t.Fail()
			}
		})
}

func TestStoreInvalidJSON(t *testing.T) {
	JsonErrorsWrapper(t, []byte(`{"password:"321cba", "maxReads": 1, "maxMinutes": 0}`), http.StatusInternalServerError,
		func(jsonResp JSONErrorResponse) {
			if jsonResp.StatusCode != http.StatusInternalServerError || !strings.Contains(jsonResp.ErrorMessage, "Syntax Error") {
				t.Logf("Expected Decode Error, got %q", jsonResp.ErrorMessage)
				t.Fail()
			}
		})
}

func TestStoreDBError(t *testing.T) {
	JsonErrorsWrapper(t, []byte(`{"password":"crash database", "maxReads": 1, "maxMinutes": 1}`), http.StatusInternalServerError,
		func(jsonResp JSONErrorResponse) {
			if jsonResp.StatusCode != http.StatusInternalServerError || !strings.Contains(jsonResp.ErrorMessage, "Database") {
				t.Logf("Expected Database Error, got %q", jsonResp.ErrorMessage)
				t.Fail()
			}
		})
}
