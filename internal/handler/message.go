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
	h.OnCommand("/sign_up", ms.SignUp)
	h.OnCommand("/login", ms.Login)
	h.OnCommand("/password", ms.Password)
	h.OnCommand("/unrecognized", ms.Password)
	h.OnCommand("/team", ms.Team)
	h.OnCommand("/create_team", ms.CreateTeam)
	h.OnCommand("/team_created", ms.TeamCreated)
	h.OnCommand("/your_team", ms.YourTeam)
	h.OnCommand("/add_user", ms.AddUser)
	h.OnCommand("/add_user_team", ms.AddUserTeam)
}

func (h *MessageHandlers) OnCommand(command string, handler model.Handler) {
	h.Handlers[command] = handler
}
