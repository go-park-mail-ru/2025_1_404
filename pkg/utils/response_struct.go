//go:generate easyjson -all

package utils

import "github.com/mailru/easyjson/jwriter"

//easyjson:json
type MessageResponse struct {
	Message string `json:"message"`
}

//easyjson:json
type ErrorResponse struct {
	Error string `json:"error"`
}

func (e ErrorResponse) MarshalEasyJSON(w *jwriter.Writer) {
	//TODO implement me
	panic("implement me")
}

//easyjson:json
type TestResponse struct {
	Message string `json:"test_message"`
}
