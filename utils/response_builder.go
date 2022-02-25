package utils

import (
	"encoding/json"
	"net/http"
)

type successResponse struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
}

type successResponseNoData struct {
	Result string `json:"result"`
}

func SendSuccessResponse(w http.ResponseWriter, result string, data interface{}) {
	responseBody := successResponse{
		Result: result,
		Data:   data,
	}

	jsonBytes, err := json.Marshal(responseBody)
	if err != nil {
		SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("content-type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		SendCustomError(w, http.StatusInternalServerError, err.Error())
	}
}

func SendSuccessResponseNoData(w http.ResponseWriter, result string) {
	var responseBody successResponseNoData
	responseBody.Result = result

	jsonBytes, err := json.Marshal(responseBody)
	if err != nil {
		SendCustomError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("content-type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		SendCustomError(w, http.StatusInternalServerError, err.Error())
	}
}
