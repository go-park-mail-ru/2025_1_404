package csrf

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateCSRF(userID string, salt string) string {
	date := time.Now().Format("2006-01-02")

	data := fmt.Sprintf("%s:%s:%s", userID, date, salt)

	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func ValidateCSRF(token string, key string, salt string) bool {
	date := time.Now().Format("2006-01-02")
	data := fmt.Sprintf("%s:%s:%s", key, date, salt)

	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(data))

	expected := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(token))
}

type CSRFResponse struct {
	CSRF string `json:"csrf_token"`
}

func GetCSRFResponse(token string) CSRFResponse  {
	return CSRFResponse{CSRF: token}
}
