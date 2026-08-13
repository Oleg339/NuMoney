package handlers

import (
	"context"
	"numoney/internal/bot"
	"numoney/internal/user"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) AddCategory(ctx context.Context, user user.User, update tgbotapi.Update) error {
	catType := "income"

	if h.states.Get(user.ID) == bot.StateAddExpenseCategory {
		catType = "expense"
	}

	_, err := h.cService.Create(ctx, user.ID, update.Message.Text, catType)
	
	if err != nil {
		return err
	}

	h.states.Set(user.ID, bot.StateNone)

	_, err = h.sendNew(update.Message.Chat.ID, bot.BaseMenu(), "Выберите действие:")

	return err
}

func (h *Handler) AddCategoryButton(_ context.Context, user user.User, update tgbotapi.Update) error {
	h.states.Set(user.ID, bot.ChoseCategoryToCreateButton)

	return h.sendEdit(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		bot.AddCategory(),
		"Выберите тип:",
	)
}

func (h *Handler) ChoseCategory(ctx context.Context, user user.User, update tgbotapi.Update) error {
	categories, err := h.cService.GetByUserID(ctx, user.ID)

	if err != nil {
		return err
	}

	if len(categories) == 0 {
		return h.sendEdit(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			bot.NoCategories(),
			"У вас нет категорий.",
		)
	}

	markup, err := bot.CategoriesJSON(categories)

	if err != nil {
		return err
	}

	_, err = h.bot.MakeRequest("editMessageReplyMarkup", tgbotapi.Params{
		"chat_id":      strconv.FormatInt(user.TelegramID, 10),
		"message_id":   strconv.Itoa(update.CallbackQuery.Message.MessageID),
		"text":         "Категории",
		"reply_markup": string(markup),
	})

	return err
}

func (h *Handler) ChoseCategoryToCreateButton(_ context.Context, user user.User, update tgbotapi.Update) error {
	if update.CallbackQuery == nil {
		return nil
	}

	h.states.Set(user.ID, update.CallbackQuery.Data)

	return h.sendEdit(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		bot.Back(),
		"Введите название категории",
	)
}
