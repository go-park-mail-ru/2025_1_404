//go:generate easyjson -all

package utils

//easyjson:json
type MessageResponse struct {
    Message string `json:"message"`
}