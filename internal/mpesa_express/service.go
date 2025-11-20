package mpesaexpress

import "context"

type Service interface {
	GetAccessToken(ctx context.Context) (error )
}

type svc struct {
	//repository
}

func NewService() Service {
	return &svc{}
}

func (s *svc) GetAccessToken(ctx context.Context) error {
	return nil
}