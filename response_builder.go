package main

import (
	"encoding/json"
	"net/http"
)

type ErrorData struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorData `json:"error"`
}

type successResponse struct {
	Result string    `json:"result"`
	Data   []Contact `json:"data"`
}

type successResponseNoData struct {
	Result string `json:"result"`
}

func sendSuccessResponse(w http.ResponseWriter, result string, data []Contact) {
	responseBody := successResponse{
		Result: result,
		Data:   data,
	}

	jsonBytes, err := json.Marshal(responseBody)
	if err != nil {
		sendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonBytes)
}

func sendSuccessResponseNoData(w http.ResponseWriter, result string) {
	var responseBody successResponseNoData
	responseBody.Result = result

	jsonBytes, err := json.Marshal(responseBody)
	if err != nil {
		sendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(jsonBytes)
}
