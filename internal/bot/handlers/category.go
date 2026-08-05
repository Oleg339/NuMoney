package handlers

import (
	"context"
	category "numoney/internal/category"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) AddCategory(ctx context.Context, update tgbotapi.Update) error {
	u, err := h.uRepo.GetByTelegramId(ctx, update.Message.From.ID)
	catType := "income"

	if h.states.Get(int64(u.ID)) == "add_expense_category" {
		catType = "expense"
	}

	err = h.cRepo.Save(ctx, &category.Category{
		Name:   update.Message.Text,
		UserId: u.ID,
		Type:   catType,
	})

	if err != nil {
		return err
	}

	return nil
}
