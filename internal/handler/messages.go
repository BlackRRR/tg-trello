package handler

import (
	"tgtrello/internal/model"
	"tgtrello/internal/service/message"
)

type MessageHandlers struct {
	Handlers map[string]model.Handler
}

func (h *MessageHandlers) GetHandler(command string) model.Handler {
	return h.Handlers[command]
}

func (h *MessageHandlers) Init(ms *message.Service) {
	h.OnCommand("/start", ms.Start)
}

func (h *MessageHandlers) OnCommand(command string, handler model.Handler) {
	h.Handlers[command] = handler
}
