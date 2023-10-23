package message

import (
	"database/sql"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tgtrello/internal/model"
)

type Service struct {
	bot  *tgbotapi.BotAPI
	rdb  *redis.Client
	repo *sql.DB
}

func NewMessageService(rdb *redis.Client, repo *sql.DB, bot *tgbotapi.BotAPI) *Service {
	return &Service{
		bot:  bot,
		rdb:  rdb,
		repo: repo,
	}
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
