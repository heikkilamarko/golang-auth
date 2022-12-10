package utils

import (
	"encoding/json"
	"net/http"
)

type DataResponse struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

func NewDataResponse(data, meta any) *DataResponse {
	return &DataResponse{data, meta}
}

func WriteResponse(w http.ResponseWriter, statusCode int, body any) {
	if body == nil {
		w.WriteHeader(statusCode)
		return
	}

	data, err := json.Marshal(body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
