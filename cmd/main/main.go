package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"tgtrello/config"
	"tgtrello/internal/handler"
	"tgtrello/internal/model"
	"tgtrello/internal/redis"
	"tgtrello/internal/repository"
	"tgtrello/internal/service/callback"
	"tgtrello/internal/service/message"
)

func main() {
	cfg := config.LoadConfig()

	logger, _ := zap.NewProduction()
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		logger.Panic("create bot instance", zap.Error(err))
	}

	db := repository.NewDB(logger, cfg)
	rdbClient, err := redis.NewClient(cfg.RedisDB.Host + ":" + cfg.RedisDB.Port)
	if err != nil {
		logger.Panic("failed to ping redis client", zap.Error(err))
	}

	logger.Info("All Databases connected successful!")

	bot.Debug = true

	logger.Info("Authorized on account", zap.String("account", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	r := handler.NewReader(logger, newMessagesHandler(message.NewMessageService(rdbClient, db, bot)), newCallbackHandler(callback.NewCallbackService(rdbClient, db, bot)))

	logger.Info("All services are running!")
	r.ReadUpdates(updates)
}

func newMessagesHandler(srv *message.Service) *handler.MessageHandlers {
	handle := handler.MessageHandlers{
		Handlers: map[string]model.Handler{},
	}

	handle.Init(srv)
	return &handle
}

func newCallbackHandler(srv *callback.Service) *handler.CallBackHandlers {
	handle := handler.CallBackHandlers{
		Handlers: map[string]model.Handler{},
	}

	handle.Init(srv)
	return &handle
}
