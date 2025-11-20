package mpesaexpress

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Service interface {
	GetAccessToken(ctx context.Context) (*AccessTokenResponse, error)
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

func (s *svc) GetAccessToken(ctx context.Context) (*AccessTokenResponse, error) {
	consumerKey := os.Getenv("CONSUMER_KEY")
	if consumerKey == "" {
		return nil, fmt.Errorf("consumer key misisng")

	}

	consumerSecret := os.Getenv("CONSUMER_SECRET")
	if consumerSecret == "" {
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
