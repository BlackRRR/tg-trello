package message

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tgtrello/internal/model"
)

type Service struct {
	bot *tgbotapi.BotAPI
}

func NewMessageService(bot *tgbotapi.BotAPI) *Service {
	return &Service{bot: bot}
}

func (m *Service) Start(s *model.Situation) error {
	msg := &tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: s.User.ID,
		},
		Text: "hello",
	}

	_, err := m.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
