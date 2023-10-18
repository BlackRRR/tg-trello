package handler

import (
	"tgtrello/internal/model"
	"tgtrello/internal/service/callback"
)

type CallBackHandlers struct {
	Handlers map[string]model.Handler
}

func (h *CallBackHandlers) GetHandler(command string) model.Handler {
	return h.Handlers[command]
}

func (h *CallBackHandlers) Init(cs *callback.Service) {
	h.OnCommand("/start", cs.Start)
	// Start commands
}

func (h *CallBackHandlers) OnCommand(command string, handler model.Handler) {
	h.Handlers[command] = handler
}
