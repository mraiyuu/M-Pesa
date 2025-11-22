package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	GetAccessToken(ctx context.Context) (*AccessTokenResponse, error)
	InitiateSTK(ctx context.Context, phoneNumber string) (*InitiateSTKResponse, error)
}

type svc struct {
	//repository
}

func NewService() Service {
	return &svc{}
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

type InitiateSTKResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

type STKPayload struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            string `json:"Amount"`
	PartyA            string `json:"PartyA"`
	PartyB            string `json:"PartyB"`
	PhoneNumber       string `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

func (s *svc) GetAccessToken(ctx context.Context) (*AccessTokenResponse, error) {
	consumerKey, ok := os.LookupEnv("CONSUMER_KEY")
	if !ok {
		return nil, fmt.Errorf("consumer key missing")
	}

	consumerSecret, ok := os.LookupEnv("CONSUMER_SECRET")
	if !ok {
		return nil, fmt.Errorf("consumer secret missing")
	}

	URL := "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, fmt.Errorf("Newrequest failed %w", err)
	}

	req.SetBasicAuth(consumerKey, consumerSecret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Default http client failed %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Read all body failed %w", err)
	}

	var tokenResp AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	return &tokenResp, nil

}

func (s *svc) InitiateSTK(ctx context.Context, phoneNumber string) (*InitiateSTKResponse, error) {
	//validate phone number before anything else
	if phoneNumber != "" {
		formattedPhone, err := formatPhoneNumber(phoneNumber)
		if err != nil {
			return nil, err
		}
		phoneNumber = formattedPhone
	}

	shortCode := os.Getenv("SHORTCODE")
	if shortCode == "" {
		return nil, fmt.Errorf("short code missing in env")
	}

	passKey := os.Getenv("PASSKEY")
	if passKey == "" {
		return nil, fmt.Errorf("passkey missing in env")
	}

	loc, _ := time.LoadLocation("Africa/Nairobi")
	timestamp := time.Now().In(loc).Format("20060102150405")

	password := base64.StdEncoding.EncodeToString([]byte(shortCode + passKey + timestamp))
	URL := "https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest"

	payload := STKPayload{
		BusinessShortCode: "174379",
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline",
		Amount:            "10",
		PartyA:            phoneNumber,
		PartyB:            "174379",
		PhoneNumber:       phoneNumber,
		CallBackURL:       "https://webhook.site/33504856-0057-4f41-a8bb-e47975b16e4e",
		AccountReference:  "Go backend Mpesa",
		TransactionDesc:   "Payment for test",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	accessToken, err := s.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to make request %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken.AccessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make http request: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	log.Printf("STK Push HTTP Status: %d", res.StatusCode)
	log.Printf("STK Push Raw Response: %s", string(body))

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stk push failed with status %d: %s", res.StatusCode, string(body))
	}

	var STKResponse InitiateSTKResponse
	if err := json.Unmarshal(body, &STKResponse); err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)

	}
	return &STKResponse, nil
}

func formatPhoneNumber(phoneNumber string) (string, error) {
	// Remove any spaces, dashes, or plus signs
	phoneNumber = strings.TrimSpace(phoneNumber)
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, "+", "")

	// If it starts with 0, replace with 254
	if strings.HasPrefix(phoneNumber, "0") {
		phoneNumber = "254" + phoneNumber[1:]
	}

	// If it doesn't start with 254, prepend it 
	if !strings.HasPrefix(phoneNumber, "254") {
		phoneNumber = "254" + phoneNumber
	}

	// Validate length (254 + 9 digits = 12 characters)
	if len(phoneNumber) != 12 {
		return "", fmt.Errorf("invalid phone number length: %d", len(phoneNumber))
	}

	// Validate all characters are digits
	if _, err := strconv.ParseInt(phoneNumber, 10, 64); err != nil {
		return "", fmt.Errorf("phone number contains non-digit characters")
	}

	return phoneNumber, nil
}
