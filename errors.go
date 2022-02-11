package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrBReq            = errors.New("bad request")
	ErrAlreadyFav      = errors.New("contact already in favorites")
	ErrAlreadyNotFav   = errors.New("contact not in favorites already")
	ErrWrongFormat     = errors.New("wrong request format, please use id={id}&action={add|remove}")
	ErrContactNotFound = errors.New("contact not found")
)

func sendCustomError(w http.ResponseWriter, status int, message string) {
	var errData ErrorData
	errData.Message = message
	var errResp ErrorResponse
	errResp.Error = errData

	jsonBytes, err := json.Marshal(errResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonBytes)
}

/*func sendInternalError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

/*func buildErrorResponse(err error) ErrorResponse {
	errMess := err.Error()
	var errData ErrorData
	errData.Message = errMess
	var errResp ErrorResponse
	errResp.Error = errMess

	return errResp
}

func sendErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	//errResp := buildErrorResponse(err)
	jsonBytes, err := json.Marshal(err.Error())
	if err != nil {
		sendInternalError(w, err)
		return
	}

	w.Write(jsonBytes)
}*/
