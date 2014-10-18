package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var (
	test_server *httptest.Server
)

type MockInsertDao struct {
}

func (my *MockInsertDao) InsertToken(token, password string, maxReads, maxMinutes int) error {
	if password == "crash database" {
		return errors.New("Fake Database Error")
	}
	return nil
}

func init() {
	// Heavily Cribbed from gotutorial.net lesson 8
	router := mux.NewRouter()
	router.Handle("/store", NewStoreHandler(new(MockInsertDao)))
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	})
	test_server = httptest.NewServer(router)
}

func TestTestServer(t *testing.T) {
	// Sanity Test. Verify fake server is up and running
	res, err := http.Get(test_server.URL)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Logf("Expected 200 got %v", res.StatusCode)
		t.Fail()
	}

}

func assertStatusCodeEquals(t *testing.T, resp *http.Response, expectedStatusCode int) {
	if resp.StatusCode != expectedStatusCode {
		t.Logf("expected status code %v got %v", expectedStatusCode, resp.StatusCode)
		t.Fail()
	}
}

func assertValidJSON(t *testing.T, resp *http.Response, jsonObj interface{}) {
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(jsonObj)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
