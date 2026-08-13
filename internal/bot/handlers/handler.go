package handlers

import (
	"context"
	"numoney/internal/bot"
	category "numoney/internal/category"
	"numoney/internal/transaction"
	"numoney/internal/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot      *tgbotapi.BotAPI
	cService category.Service
	tService transaction.Service
	states   *bot.State
}

func NewHandler(bot *tgbotapi.BotAPI, cService category.Service, tService transaction.Service, states *bot.State) *Handler {
	return &Handler{bot: bot, cService: cService, tService: tService, states: states}
}

func (h *Handler) sendEdit(chatID int64, messageID int, keyboard tgbotapi.InlineKeyboardMarkup, text string) error {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ParseMode = "HTML"

	if len(keyboard.InlineKeyboard) > 0 {
		msg.ReplyMarkup = &keyboard
	}

	_, err := h.bot.Send(msg)

	return err
}

func (h *Handler) sendNew(chatID int64, keyboard tgbotapi.InlineKeyboardMarkup, text string) (tgbotapi.Message, error) {
	msgNew := tgbotapi.NewMessage(chatID, text)

	if len(keyboard.InlineKeyboard) > 0 {
		msgNew.ReplyMarkup = keyboard
	}

	return h.bot.Send(msgNew)
}

func (h *Handler) InitMessage(_ context.Context, user user.User, update tgbotapi.Update) error {
	h.states.Set(user.ID, bot.StateNone)

	if update.Message == nil {
		err := h.sendEdit(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			bot.BaseMenu(),
			"Выберите действие:",
		)

		return err
	}

	_, err := h.sendNew(update.Message.Chat.ID, bot.BaseMenu(), "Выберите действие:")

	return err

}
