package callback

import (
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"tgtrello/internal/model"
	"tgtrello/internal/repository"
)

type Service struct {
	log   *zap.Logger
	texts map[string]string
	bot   *tgbotapi.BotAPI
	rdb   *redis.Client
	repo  *repository.PGRepository
}

func NewCallbackService(log *zap.Logger, rdb *redis.Client, repo *repository.PGRepository, bot *tgbotapi.BotAPI, texts map[string]string) *Service {
	return &Service{
		log:   log,
		rdb:   rdb,
		repo:  repo,
		bot:   bot,
		texts: texts,
	}
}

func (c *Service) Start(s *model.Situation) error {
	return nil
}
