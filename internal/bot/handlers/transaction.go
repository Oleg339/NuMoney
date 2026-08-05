package handlers

import (
	"context"
	"numoney/internal/bot"
	"numoney/internal/transaction"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) AddTransactionAmount(ctx context.Context, update tgbotapi.Update) error {
	var keyboard tgbotapi.InlineKeyboardMarkup
	var msg tgbotapi.EditMessageTextConfig

	text := strings.Replace(update.Message.Text, ",", ".", 1)
	amount, err := strconv.ParseFloat(text, 64)

	u, err := h.uRepo.GetByTelegramId(ctx, update.Message.From.ID)
	userId := u.ID

	if err != nil {
		text = "Введите корректное число:"

		h.states.Set(userId, bot.StateAddTransactionAmount)
		keyboard = tgbotapi.InlineKeyboardMarkup{}

	} else {
		err := h.tRepo.Save(ctx, &transaction.Transaction{
			Amount:     int(amount * 100),
			CategoryId: h.states.GetCategory(userId),
			UserId:     u.ID,
		})

		if err != nil {
			return err
		}

		h.states.SetCategory(userId, 0)
		h.states.Set(userId, bot.StateNone)
	}

	if len(keyboard.InlineKeyboard) > 0 {
		msg.ReplyMarkup = &keyboard
	}

	h.bot.Send(msg)

	return nil
}


