//go:generate easyjson -all

package utils

//easyjson:json
type ErrorResponse struct {
    Error string `json:"error"`
}