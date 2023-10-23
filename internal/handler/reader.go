package handler

import (
	"database/sql"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"tgtrello/internal/model"
	"tgtrello/internal/service/callback"
	"tgtrello/internal/service/message"
)

type Reader struct {
	texts    map[string]string
	logger   *zap.Logger
	msg      *MessageHandlers
	callback *CallBackHandlers
}

func NewReader(log *zap.Logger, rdb *redis.Client, repo *sql.DB, bot *tgbotapi.BotAPI, texts map[string]string) *Reader {
	return &Reader{
		logger:   log,
		msg:      newMessagesHandler(message.NewMessageService(rdb, repo, bot, texts)),
		callback: newCallbackHandler(callback.NewCallbackService(rdb, repo, bot, texts)),
	}
}

func (r *Reader) ReadUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		go r.updateActions(update)
	}
}

func (r *Reader) updateActions(update tgbotapi.Update) {
	if update.Message != nil {
		s := setMessageSituation(update.Message)

		handler := r.msg.GetHandler(update.Message.Text)
		err := handler(s)
		if err != nil {
			r.logger.Error("failed to get handler", zap.Error(err))
		}
	}

	if update.CallbackQuery != nil {
		s := setCallbackSituation(update.CallbackQuery)

		handler := r.callback.GetHandler(update.Message.Text)
		err := handler(s)
		if err != nil {
			r.logger.Error("failed to get handler", zap.Error(err))
		}
	}
}

func setMessageSituation(message *tgbotapi.Message) *model.Situation {
	return &model.Situation{
		Message: message,
		User:    &model.User{ID: message.Chat.ID},
	}
}

func setCallbackSituation(callback *tgbotapi.CallbackQuery) *model.Situation {
	return &model.Situation{
		CallbackQuery: callback,
		User:          &model.User{ID: callback.Message.Chat.ID},
	}
}

func newMessagesHandler(srv *message.Service) *MessageHandlers {
	handle := MessageHandlers{
		Handlers: map[string]model.Handler{},
	}

	handle.Init(srv)
	return &handle
}

func newCallbackHandler(srv *callback.Service) *CallBackHandlers {
	handle := CallBackHandlers{
		Handlers: map[string]model.Handler{},
	}

	handle.Init(srv)
	return &handle
}
