package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

type errorData struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorData `json:"error"`
}

var (
	ErrAlreadyFav        = errors.New("contact already in favorites")
	ErrAlreadyNotFav     = errors.New("contact not in favorites already")
	ErrFavWrongFormat    = errors.New("wrong request format, please use id={id}&action=add|remove")
	ErrFilterWrongFormat = errors.New("wrong request format, please use field=name|phone&value={string}")
	ErrContactNotFound   = errors.New("contact not found")
)

func SendCustomError(w http.ResponseWriter, status int, message string) {
	var errData errorData
	errData.Message = message
	var errResp errorResponse
	errResp.Error = errData

	jsonBytes, err := json.Marshal(errResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	_, err = w.Write(jsonBytes)
	if err != nil {
		SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
