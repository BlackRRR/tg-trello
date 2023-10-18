package callback

import "tgtrello/internal/model"

type Service struct{}

func NewCallbackService() *Service {
	return &Service{}
}

func (c *Service) Start(s *model.Situation) error {
	return nil
}
