package jsonw

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Message is interface for valid JSON struct
type Message interface {
	HTTPStatus() int
}

// SuccessJSON struct contains message json format for success/fail
// For more information check:
type SuccessJSON struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Code   int         `json:"-"`
}

// HTTPStatus return http status for success response (200)
func (s SuccessJSON) HTTPStatus() int {
	return http.StatusOK
}

// FailJSON struct is same what success, but with correct name in function
type FailJSON SuccessJSON

// HTTPStatus return hhtp status for fail response (200)
func (f FailJSON) HTTPStatus() int {
	return http.StatusUnprocessableEntity
}

// ErrorJSON struct contains message json format for error
// For more information check:
type ErrorJSON struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message interface{} `json:"message"`
	Code    int         `json:"-"`
}

// HTTPStatus return http status for error response - code value
func (s ErrorJSON) HTTPStatus() int {
	return s.Code
}

// Success create json message with success status and data
func Success(w http.ResponseWriter, data interface{}, code int) {
	var json = SuccessJSON{
		Status: "success",
		Data:   data,
		Code:   code,
	}
	RespondWithJSON(w, json)
}

// Fail create json message with faile status and data
func Fail(w http.ResponseWriter, data interface{}) {
	var json = FailJSON{
		Status: "fail",
		Data:   data,
	}
	RespondWithJSON(w, json)
}

// FailValidation create json message with fail status and data
// contains validation with param data
func FailValidation(w http.ResponseWriter, data interface{}) {
	var json = FailJSON{
		Status: "fail",
		Data: map[string]interface{}{
			"validation": data,
		},
	}
	RespondWithJSON(w, json)
}

// Error create json message with erorr. Optional params are data and code. If
// error doesn't contains data or/and code just set nil
func Error(w http.ResponseWriter, message interface{}, data interface{}, code int) {
	var json = ErrorJSON{
		Status:  "error",
		Message: message,
		Data:    data,
		Code:    code,
	}
	RespondWithJSON(w, json)
}

// RespondWithJSON sends value v as a JSON response.
func RespondWithJSON(w http.ResponseWriter, v Message) {
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		log.Print("failed to encode response as JSON: ", err)
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(v.HTTPStatus())
	buf.WriteTo(w)
}
