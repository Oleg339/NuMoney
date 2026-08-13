package handlers

import (
	"context"
	"errors"
	"numoney/internal/bot"
	"numoney/internal/category"
	"numoney/internal/user"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) AddTransactionButton(ctx context.Context, user user.User, update tgbotapi.Update) error {
	h.states.Set(user.ID, bot.StateChoseTransactionCategory)

	kind := category.TypeExpense

	if update.CallbackQuery.Data == bot.AddIncomeButton {
		kind = category.TypeIncome
	}

	categories, err := h.cService.GetByUserIDAndType(ctx, user.ID, kind)

	if err != nil {
		return err
	}

	markup, err := bot.ChoseCategoryJSON(categories)

	if err != nil {
		return err
	}

	_, err = h.bot.MakeRequest("editMessageReplyMarkup", tgbotapi.Params{
		"chat_id":      strconv.FormatInt(user.TelegramID, 10),
		"message_id":   strconv.Itoa(update.CallbackQuery.Message.MessageID),
		"text":         "Выберите категорию",
		"reply_markup": string(markup),
	})

	return err
}

func (h *Handler) StateChoseTransactionCategory(ctx context.Context, user user.User, update tgbotapi.Update) error {
	if update.CallbackQuery == nil {
		return nil
	}

	categoryID, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return h.InitMessage(ctx, user, update)
	}

	cat, err := h.cService.Find(ctx, int64(categoryID), user.ID)
	if err != nil {
		_, sendErr := h.sendNew(update.CallbackQuery.Message.Chat.ID, tgbotapi.InlineKeyboardMarkup{}, "Категория не найдена")
		if sendErr != nil {
			return sendErr
		}

		return h.InitMessage(ctx, user, update)
	}

	h.states.SetCategory(user.ID, cat.ID)
	h.states.Set(user.ID, bot.StateAddTransactionAmount)

	return h.sendEdit(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		tgbotapi.InlineKeyboardMarkup{},
		"Введите сумму:",
	)
}

func (h *Handler) AddTransactionAmount(ctx context.Context, user user.User, update tgbotapi.Update) error {
	text := strings.Replace(update.Message.Text, ",", ".", 1)
	amount, err := strconv.ParseFloat(text, 64)

	if err != nil {
		_, err := h.sendNew(
			update.Message.Chat.ID,
			tgbotapi.InlineKeyboardMarkup{},
			"Введите корректное число:",
		)

		return err
	}

	catID := h.states.GetCategory(user.ID)

	_, err = h.tService.Create(ctx, amount, catID, user.ID)

	switch {
	case err == nil:
		h.states.Set(user.ID, bot.StateNone)

		h.states.SetCategory(user.ID, 0)

		return h.InitMessage(ctx, user, update)

	case errors.Is(err, category.ErrNotFound):
		_, err := h.sendNew(
			update.Message.Chat.ID,
			tgbotapi.InlineKeyboardMarkup{},
			"Категории не существует",
		)

		h.states.Set(user.ID, bot.StateNone)

		h.states.SetCategory(user.ID, 0)

		if err != nil {
			return err
		}
		
		return nil 
	}

	h.states.Set(user.ID, bot.StateNone)

	h.states.SetCategory(user.ID, 0)
	return err
}
