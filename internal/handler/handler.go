package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"tgtrello/internal/model"
)

type Reader struct {
	logger   *zap.Logger
	msg      *MessageHandlers
	callback *CallBackHandlers
}

func NewReader(log *zap.Logger, m *MessageHandlers, c *CallBackHandlers) *Reader {
	return &Reader{
		logger:   log,
		msg:      m,
		callback: c,
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
