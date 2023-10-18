package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler func(situation *Situation) error

type Situation struct {
	Message       *tgbotapi.Message       `json:"message,omitempty"`
	CallbackQuery *tgbotapi.CallbackQuery `json:"callback_query,omitempty"`
	User          *User                   `json:"user,omitempty"`
	Command       string                  `json:"command,omitempty"`
	//Params        *Parameters             `json:"params,omitempty"`
	Err error `json:"err,omitempty"`
}
