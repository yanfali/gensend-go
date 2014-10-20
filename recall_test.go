package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func getRecallUrl() string {
	return test_server.URL + getBaseAPIUrl() + "/recall"
}

func TestRecallOK(t *testing.T) {
	resp, err := http.Get(getRecallUrl() + "/tokenabc")
	if resp.StatusCode != http.StatusOK {
		t.Log(err)
		t.Fail()
	}
}

func TestRecallMaxReadsReached(t *testing.T) {
	// Read token once
	resp, err := http.Get(getRecallUrl() + "/tokenReadOnce")
	if resp.StatusCode != http.StatusOK {
		t.Log(err)
		t.Fail()
	}

	d := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	result := GensendgoRow{}
	if err := d.Decode(&result); err != nil {
		t.Log(err)
		t.Fail()
	}
	if result.MaxReads != 0 {
		t.Log("expected 0 reads left, got %v", result)
		t.Fail()
	}

	// Second read should fail with no token returned
	resp, err = http.Get(getRecallUrl() + "/tokenReadOnce")
	if resp.StatusCode != http.StatusNotFound {
		t.Log(err)
		t.Fail()
	}
}

func TestRecallTokenNotFound(t *testing.T) {
	resp, err := http.Get(getRecallUrl() + "/tokenBogus")
	if resp.StatusCode != http.StatusNotFound {
		t.Log(err)
		t.Fail()
	}
}
