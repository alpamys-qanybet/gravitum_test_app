package model

import (
	"net/http"
	"strings"
)

type Response struct {
	StatusCode int          `json:"status_code"`
	StatusText string       `json:"status_text"`
	Data       *interface{} `json:"data,omitempty"`
}

func WrapResponse(statusCode int, data interface{}) Response {
	statusText := http.StatusText(statusCode)
	if statusText == "" {
		statusText = "Unknown Status"
	}

	statusText = strings.ToLower(statusText)
	statusText = strings.ReplaceAll(statusText, " ", "_")

	if data != nil {
		return Response{
			StatusCode: statusCode,
			StatusText: statusText,
			Data:       &data,
		}
	}

	return Response{
		StatusCode: statusCode,
		StatusText: statusText,
	}
}
