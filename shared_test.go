package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

var (
	test_server       *httptest.Server
	mockAccountingDao *MockAccountingDao
)

const (
	POST = "POST"
	GET  = "GET"
)

type MockInsertDao struct {
}

func (my *MockInsertDao) InsertToken(token, password string, maxReads, maxMinutes int) error {
	if password == "crash database" {
		return errors.New("Fake Database Error")
	}
	return nil
}

type MockAccountingDao struct {
	rows map[string]*GensendgoRow
}

func getUTCTime() time.Time {
	return time.Now().UTC()
}

func NewMockAccountingDao() *MockAccountingDao {
	dao := new(MockAccountingDao)
	dao.rows = make(map[string]*GensendgoRow)
	dao.rows["tokenabc"] = &GensendgoRow{"tokenabc", 1, 1, getUTCTime(), getUTCTime(), "password"}
	dao.rows["tokenReadOnce"] = &GensendgoRow{"tokenReadOnce", 1, 1, getUTCTime(), getUTCTime().Add(time.Minute * 1), "password"}
	return dao
}

func (my *MockAccountingDao) UpdateMaxReadCount(aRow *GensendgoRow) error {
	return nil
}
func (my *MockAccountingDao) DeleteById(id string) error {
	my.rows[id].MaxReads--
	log.Printf("DeleteById: %q", id)
	return nil
}
func (my *MockAccountingDao) FetchValidRowsById(token string) ([]GensendgoRow, error) {
	return []GensendgoRow{}, nil
}

func (my *MockAccountingDao) FetchValidRowById(token string) ([]GensendgoRow, error) {
	rows := []GensendgoRow{}
	if row, ok := my.rows[token]; ok {
		if row.MaxReads > 0 {
			rows = append(rows, *row)
		}
	}
	return rows, nil
}

func init() {
	// Heavily Cribbed from gotutorial.net lesson 8
	router := mux.NewRouter()
	router.Handle(getBaseAPIUrl()+"/store", appHandler{NewStoreHandler(new(MockInsertDao))})
	mockAccountingDao = NewMockAccountingDao()
	router.Handle(getBaseAPIUrl()+"/recall/{token}", appHandler{NewRecallHandler(mockAccountingDao)})
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

func getClientResponse(t *testing.T, verb string, content []byte) *http.Response {
	// Boiler Plate for Setting up JSON Request
	req, err := http.NewRequest(verb, getStoreUrl(), bytes.NewBuffer(content))
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
