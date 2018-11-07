package response

import "net/http"

// ErrorResult is for creating a JSON response of `{ "statusCode": 400, "errorCode": "reason_failed", "message": "didn't work" }`.
type ErrorResult struct {
	StatusCode int                    `json:"statusCode"`
	ErrorCode  string                 `json:"errorCode"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// ErrorResponse matches the standard:
// {
// 		"error": {
// 		"statusCode": 404,
// 		"errorCode": "UserNotFound",
// 		"message": "The account could not be found"
// 	}
// }
//
type ErrorResponse struct {
	Error ErrorResult `json:"error"`
}

// NewErrorResponse creates a standard error response.
func NewErrorResponse(statusCode int, errorCode string, message string) ErrorResponse {
	return NewErrorResponseWithData(statusCode, errorCode, message, nil)
}

// NewErrorResponseWithData creates a standard error response.
func NewErrorResponseWithData(statusCode int, errorCode string, message string, data map[string]interface{}) ErrorResponse {
	return ErrorResponse{
		Error: ErrorResult{
			StatusCode: statusCode,
			ErrorCode:  errorCode,
			Message:    message,
			Data:       data,
		},
	}
}

// Error returns a JSON ErrorResult.
func Error(code string, err error, w http.ResponseWriter, status int) error {
	return ErrorWithData(code, err, nil, w, status)
}

// ErrorWithData returns a JSON ErrorResult.
func ErrorWithData(code string, err error, data map[string]interface{}, w http.ResponseWriter, status int) error {
	response := NewErrorResponseWithData(status, code, err.Error(), data)
	return JSON(response, w, status)
}

// ErrorString returns a JSON ErrorResult.
func ErrorString(code string, message string, w http.ResponseWriter, status int) error {
	return ErrorStringWithData(code, message, nil, w, status)
}

// ErrorStringWithData returns a JSON ErrorResult.
func ErrorStringWithData(code string, message string, data map[string]interface{}, w http.ResponseWriter, status int) error {
	response := NewErrorResponseWithData(status, code, message, data)
	return JSON(response, w, status)
}
