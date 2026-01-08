package natdemo

import (
	"context"
)

type NATSService struct {
}

func NewNATSService() *NATSService {
	return &NATSService{}
}

func (s *NATSService) Create(ctx context.Context) error {
	return nil
}
