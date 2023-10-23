package callback

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

func NewCallbackService(rdb *redis.Client, repo *sql.DB, bot *tgbotapi.BotAPI, texts map[string]string) *Service {
	return &Service{
		rdb:   rdb,
		repo:  repo,
		bot:   bot,
		texts: texts,
	}
}

func (c *Service) Start(s *model.Situation) error {
	return nil
}
