package bot

import (
	"encoding/json"
	"fmt"
	"numoney/internal/category"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Keyboards struct {
}

func BaseMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", "Statistics"),
			tgbotapi.NewInlineKeyboardButtonData("Категории", "Categories"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔴 Расход", "Expense"),
			tgbotapi.NewInlineKeyboardButtonData("🟢 Доход", "Income"),
		),
	)
}

func AddCategory() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔴 Расход", "add_expense_category"),
			tgbotapi.NewInlineKeyboardButtonData("🟢 Доход", "add_income_category"),
		),
	)
}

func ChoseStatPeriod() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("За сегодня", "today_stats"),
			tgbotapi.NewInlineKeyboardButtonData("За неделю", "week_stats"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("За месяц", "month_stats"),
			tgbotapi.NewInlineKeyboardButtonData("За всё время", "all_stats"),
		),
	)
}

func Back() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back"),
		),
	)
}

func NoCategories() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить", "add_category"),
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back"),
		),
	)
}

func CategoriesJSON(c []category.Category) ([]byte, error) {
	type inlineBtn struct {
		Text         string `json:"text"`
		CallbackData string `json:"callback_data,omitempty"`
		Style        string `json:"style,omitempty"`
	}

	row1 := []inlineBtn{
		{Text: "⬅️ Назад", CallbackData: "back"},
		{Text: "➕ Добавить", CallbackData: "add_category"},
	}

	var rows [][]inlineBtn
	var currentRow []inlineBtn

	for _, category := range c {
		b := inlineBtn{
			Text:         category.Name,
			CallbackData: fmt.Sprint(category.ID),
		}
		if category.Type == "income" {
			b.Style = "success"
		} else {
			b.Style = "danger"
		}

		currentRow = append(currentRow, b)

		if len(currentRow) == 3 {
			rows = append(rows, currentRow)
			currentRow = nil
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	allRows := append([][]inlineBtn{row1}, rows...)

	markup, err := json.Marshal(map[string]any{
		"inline_keyboard": allRows,
	})

	if err != nil {
		return nil, err
	}

	return markup, nil
}

func ChoseCategoryJSON(c []category.Category) ([]byte, error) {
	type inlineBtn struct {
		Text         string `json:"text"`
		CallbackData string `json:"callback_data,omitempty"`
		Style        string `json:"style,omitempty"`
	}

	row1 := []inlineBtn{
		{Text: "⬅️ Назад", CallbackData: "back"},
		// {Text: "➕ Добавить", CallbackData: "add_category"},
	}

	var rows [][]inlineBtn
	var currentRow []inlineBtn

	for _, category := range c {
		b := inlineBtn{
			Text:         category.Name,
			CallbackData: fmt.Sprint(category.ID),
		}
		if category.Type == "income" {
			b.Style = "success"
		} else {
			b.Style = "danger"
		}

		currentRow = append(currentRow, b)

		if len(currentRow) == 3 {
			rows = append(rows, currentRow)
			currentRow = nil
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	allRows := append([][]inlineBtn{row1}, rows...)

	markup, err := json.Marshal(map[string]any{
		"inline_keyboard": allRows,
	})

	if err != nil {
		return nil, err
	}

	return markup, nil
}
