package message

import (
	"database/sql"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tgtrello/internal/model"
)

type Service struct {
	texts map[string]string
	bot   *tgbotapi.BotAPI
	rdb   *redis.Client
	repo  *sql.DB
}

func NewMessageService(rdb *redis.Client, repo *sql.DB, bot *tgbotapi.BotAPI, texts map[string]string) *Service {
	return &Service{
		bot:   bot,
		rdb:   rdb,
		repo:  repo,
		texts: texts,
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
