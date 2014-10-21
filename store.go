package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/unrolled/render.v1"
)

// Store the db handle within the handler struct associated with the handler
// No need to pass around a global
type StoreHandler struct {
	dao          InsertDao
	urlGenerator *UrlGenerator
}

type JSONStoreRequest struct {
	Password   string `json:"password"`
	MaxReads   int    `json:"maxReads"`
	MaxMinutes int    `json:"maxMinutes"`
}

type JSONStoreResponse struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

type JSONErrorResponse struct {
	StatusCode   int    `json:"statusCode"`
	ErrorMessage string `json:"errorMessage"`
}

// Create a New Recall Handler with the appropriate db handle
func NewStoreHandler(dao InsertDao) *StoreHandler {
	return &StoreHandler{dao, new(UrlGenerator)}
}

func validateJsonRequest(jsonRequest *JSONStoreRequest) (err error, invalid bool) {
	if jsonRequest.MaxMinutes <= 0 {
		return errors.New(fmt.Sprintf("maxMinutes of %d should be greater than 0", jsonRequest.MaxMinutes)), true
	}
	if jsonRequest.MaxReads <= 0 {
		return errors.New(fmt.Sprintf("maxReads of %d should be greater than 0", jsonRequest.MaxReads)), true
	}

	if len(jsonRequest.Password) == 0 {
		return errors.New(fmt.Sprintf("Empty Password!")), true
	}
	return nil, false
}

// HTTP Handler
func (my *StoreHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	addCORS(w)
	decoder := json.NewDecoder(req.Body)
	var jsonRequest JSONStoreRequest
	r := render.New(render.Options{})
	err := decoder.Decode(&jsonRequest)
	if err != nil {
		errMsg := fmt.Sprintf("500 Internal Server Error: Decode Error (%v)", err)
		log.Printf("%#v", err)
		if serr, ok := err.(*json.UnmarshalTypeError); ok {
			errMsg = fmt.Sprintf("Error decoding field expected type %v got %s", serr.Type, serr.Value)
		}
		log.Printf("%s", errMsg)
		r.JSON(w, http.StatusInternalServerError, JSONErrorResponse{http.StatusInternalServerError, errMsg})
		return
	}

	if err, invalid := validateJsonRequest(&jsonRequest); invalid {
		r.JSON(w, http.StatusBadRequest, JSONErrorResponse{http.StatusBadRequest, fmt.Sprintf("%s", err)})
		return
	}

	token := my.urlGenerator.Generate(jsonRequest.Password)
	log.Printf("%v", jsonRequest)
	err = my.dao.InsertToken(token, jsonRequest.Password, jsonRequest.MaxReads, jsonRequest.MaxMinutes)
	if err != nil {
		errMsg := fmt.Sprintf("Store Password Error (%v)", err)
		log.Printf("%s", errMsg)
		r.JSON(w, http.StatusInternalServerError, JSONErrorResponse{http.StatusInternalServerError, errMsg})
		return
	}
	var results JSONStoreResponse = JSONStoreResponse{token, fmt.Sprintf("/recall/%s", token)}
	r.JSON(w, http.StatusOK, results)
}
