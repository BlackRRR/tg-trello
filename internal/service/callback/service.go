package callback

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

func NewCallbackService(rdb *redis.Client, repo *sql.DB, bot *tgbotapi.BotAPI) *Service {
	return &Service{
		rdb:  rdb,
		repo: repo,
		bot:  bot,
	}
}

func (c *Service) Start(s *model.Situation) error {
	return nil
}
