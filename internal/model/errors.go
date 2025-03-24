package model

import (
	"gravitum-test-app/pkg/errors"
	"net/http"
	"strings"
)

var (
	StatusUnauthorized                   string = "unauthorized"
	ErrSecurityUnauthorized              error  = errors.New("err.security.unauthorized")
	ErrSecurityUnauthorizedNoHeader      error  = errors.New("err.security.unauthorized-no-header")
	ErrSecurityUnauthorizedInvalidHeader error  = errors.New("err.security.unauthorized-invalid-header")
	ErrSecurityAbsentSecret              error  = errors.New("err.security.absent-secret")
	ErrSecurityInvalidSecret             error  = errors.New("err.security.invalid-secret")
	ErrResponseUnexpectedStatusCode      error  = errors.New("err.response.unexpected_status_code")
	ErrRequestInvalidUrlParams           error  = errors.New("err.request.invalid_url_params")
	ErrRequestInvalidBodyParams          error  = errors.New("err.request.invalid_body_params")
	ErrRequestNameRequired               error  = errors.New("err.request.name_required")
	ErrNoUserWithSuchId                  error  = errors.New("err.user.no_user_with_such_id")
	ErrSqlNoRows                         error  = errors.New("err.sql.no_rows")
)

type ErrorResponse struct {
	Err *Errors `json:"errors"`
}

type Errors struct {
	StatusCode int     `json:"status_code"`
	StatusText string  `json:"status_text"`
	Err        *string `json:"err,omitempty"`
}

func WrapError(statusCode int, err string) ErrorResponse {
	statusText := http.StatusText(statusCode)
	if statusText == "" {
		statusText = "Unknown Status"
	}

	statusText = strings.ToLower(statusText)
	statusText = strings.ReplaceAll(statusText, " ", "_")

	return ErrorResponse{
		Err: &Errors{
			StatusCode: statusCode,
			StatusText: statusText,
			Err:        &err,
		},
	}
}
