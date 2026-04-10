package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type PaymentClient struct {
	baseURL string
	client  *http.Client
}

type CreatePaymentRequest struct {
	OrderID string `json:"order_id"`
	Amount  int64  `json:"amount"`
}

type CreatePaymentResponse struct {
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
}

func NewPaymentClient(baseURL string) *PaymentClient {
	return &PaymentClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 2 * time.Second, 
		},
	}
}

func (c *PaymentClient) CreatePayment(orderID string, amount int64) (string, error) {
	reqBody := CreatePaymentRequest{
		OrderID: orderID,
		Amount:  amount,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := c.client.Post(
		c.baseURL+"/payments",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("payment service error")
	}

	var paymentResp CreatePaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return "", err
	}

	return paymentResp.Status, nil
}
