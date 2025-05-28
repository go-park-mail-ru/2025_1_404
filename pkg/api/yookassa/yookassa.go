package yookassa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_404/config"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strconv"
)

type yookassaRepo struct {
	secret string
	shopId string
}

func New(cfg *config.YookassaConfig) *yookassaRepo {
	return &yookassaRepo{secret: cfg.Secret, shopId: cfg.ShopId}
}

func (repo *yookassaRepo) CreatePayment(amount int, description string, redirectUri string) (*CreatePaymentResponse, error) {
	client := &http.Client{}
	jsonData, _ := json.Marshal(&CreatePaymentRequest{
		Amount: Amount{
			Currency: "RUB",
			Value:    strconv.Itoa(amount),
		},
		Capture: true,
		Confirmation: Confirmation{
			Type:      "redirect",
			ReturnUri: redirectUri,
		},
		Description: description,
	})
	req, err := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(repo.shopId, repo.secret)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Idempotence-Key", uuid.New().String())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create payment uri: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	paymentResponse := &CreatePaymentResponse{}
	err = json.Unmarshal(body, paymentResponse)
	if err != nil {
		return nil, err
	}
	return paymentResponse, nil
}

func (repo *yookassaRepo) GetPayment(paymentId string) (*CreatePaymentResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.yookassa.ru/v3/payments/%s", paymentId), nil)
	req.SetBasicAuth(repo.shopId, repo.secret)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create payment uri: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	paymentResponse := &CreatePaymentResponse{}
	err = json.Unmarshal(body, paymentResponse)
	if err != nil {
		return nil, err
	}
	return paymentResponse, nil
}
