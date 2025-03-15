package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrCouldNotReadBody  = errors.New("could not read body")
	ErrCouldNotParseBody = errors.New("could not parse body")
)

type httpResp struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func sendResponse(rw http.ResponseWriter, status int, data interface{}, message string) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	out, err := json.Marshal(httpResp{Status: status, Data: data, Message: message})
	if err != nil {
		sendErrorResponse(rw, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	rw.Write(out)
}
func sendErrorResponse(rw http.ResponseWriter, status int, data interface{}, message string) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	out, _ := json.Marshal(httpResp{Status: status,
		Message: message,
		Data:    data})

	rw.Write(out)
}

func getBodyWithType[T any](r *http.Request) (T, error) {
	var v T
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return v, newError(http.StatusBadRequest, ErrCouldNotReadBody.Error())
	}
	err = json.Unmarshal(body, &v)
	if err != nil {
		return v, newError(http.StatusBadRequest, ErrCouldNotParseBody.Error())
	}
	return v, nil
}

func newError(status int, message string) error {
	return fmt.Errorf("%d %s", status, message)
}
