package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (my *StoreHandler) handler(w http.ResponseWriter, req *http.Request) (int, error) {
	addCORS(w)
	decoder := json.NewDecoder(req.Body)
	var jsonRequest JSONStoreRequest
	r := render.New(render.Options{})
	if err := decoder.Decode(&jsonRequest); err != nil {
		errMsg := fmt.Sprintf("%#v", err)
		switch err.(type) {
		case *json.SyntaxError:
			errMsg = fmt.Sprintf("Syntax Error %v", err)
		case *json.UnmarshalTypeError:
			serr := err.(*json.UnmarshalTypeError)
			errMsg = fmt.Sprintf("Type Error expected type %v got %q", serr.Type, serr.Value)
		}
		return http.StatusInternalServerError, errors.New(errMsg)
	}
	if err, invalid := validateJsonRequest(&jsonRequest); invalid {
		return http.StatusBadRequest, err
	}
	token := my.urlGenerator.Generate(jsonRequest.Password)
	if err := my.dao.InsertToken(token, jsonRequest.Password, jsonRequest.MaxReads, jsonRequest.MaxMinutes); err != nil {
		return http.StatusInternalServerError, err
	}
	r.JSON(w, http.StatusOK, JSONStoreResponse{token, fmt.Sprintf("/recall/%s", token)})
	return http.StatusOK, nil
}
