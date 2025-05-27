//go:generate easyjson -all

package utils

//easyjson:json
type MessageResponse struct {
    Message string `json:"message"`
}

//easyjson:json
type ErrorResponse struct {
    Error string `json:"error"`
}

//easyjson:json
type TestResponse struct {
	Message string `json:"test_message"`
}