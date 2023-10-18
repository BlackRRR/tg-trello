package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"tgtrello/config"
	"tgtrello/internal/handler"
	"tgtrello/internal/model"
	"tgtrello/internal/service/callback"
	"tgtrello/internal/service/message"
)

func main() {
	cfg := config.LoadConfig()

	logger, _ := zap.NewProduction()
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	r := handler.NewReader(logger, newMessagesHandler(message.NewMessageService(bot)), newCallbackHandler(callback.NewCallbackService()))
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
