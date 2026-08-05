package handlers

import (
	category "numoney/internal/category"
	"numoney/internal/bot"
	"numoney/internal/transaction"
	"numoney/internal/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot    *tgbotapi.BotAPI
	uRepo  user.Repository
	cRepo  category.Repository
	tRepo  transaction.Repository
	states *bot.State
}

func NewHandler(bot *tgbotapi.BotAPI, uRepo user.Repository, cRepo category.Repository, tRepo transaction.Repository, states *bot.State) Handler {
	return Handler{bot, uRepo, cRepo, tRepo, states}
}