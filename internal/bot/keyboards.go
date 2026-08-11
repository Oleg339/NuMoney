package bot

import (
	"encoding/json"
	"numoney/internal/category"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	AddExpenseButton            = "add_expense_button"
	AddIncomeButton             = "add_income_button"
	ChoseCategoryToCreateButton = "chose_category_to_create_button"
	ChoseStatisticsButton       = "chose_statistics_button"
	ChoseCategoryButton         = "chose_category_button"
	AddCategoryButton           = "add_category"
	TodayStatsButton            = "today_stats"
	WeekStatsButton             = "week_stats"
	MonthStatsButton            = "month_stats"
	AllStatsButton              = "all_stats"
)

type inlineBtn struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	Style        string `json:"style,omitempty"`
}

func BaseMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", ChoseStatisticsButton),
			tgbotapi.NewInlineKeyboardButtonData("Категории", ChoseCategoryButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔴 Расход", AddExpenseButton),
			tgbotapi.NewInlineKeyboardButtonData("🟢 Доход", AddIncomeButton),
		),
	)
}

func AddCategory() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", BackButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔴 Расход", StateAddExpenseCategory),
			tgbotapi.NewInlineKeyboardButtonData("🟢 Доход", StateAddIncomeCategory),
		),
	)
}

func ChoseStatPeriod() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", BackButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("За сегодня", TodayStatsButton),
			tgbotapi.NewInlineKeyboardButtonData("За неделю", WeekStatsButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("За месяц", MonthStatsButton),
			tgbotapi.NewInlineKeyboardButtonData("За всё время", AllStatsButton),
		),
	)
}

func Back() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", BackButton),
		),
	)
}

func NoCategories() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить", AddCategoryButton),
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", BackButton),
		),
	)
}

func buildCategoryKeyboard(firstRow []inlineBtn, categories []category.Category) ([]byte, error) {
	var rows [][]inlineBtn
	var currentRow []inlineBtn

	for _, cat := range categories {
		b := inlineBtn{
			Text:         cat.Name,
			CallbackData: strconv.Itoa(cat.ID),
		}
		if cat.Type == "income" {
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

	return json.Marshal(map[string]any{
		"inline_keyboard": append([][]inlineBtn{firstRow}, rows...),
	})
}

func CategoriesJSON(c []category.Category) ([]byte, error) {
	return buildCategoryKeyboard([]inlineBtn{
		{Text: "⬅️ Назад", CallbackData: BackButton},
		{Text: "➕ Добавить", CallbackData: AddCategoryButton},
	}, c)
}

func ChoseCategoryJSON(c []category.Category) ([]byte, error) {
	return buildCategoryKeyboard([]inlineBtn{
		{Text: "⬅️ Назад", CallbackData: BackButton},
	}, c)
}
