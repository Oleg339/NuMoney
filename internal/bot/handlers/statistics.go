package handlers

import (
	"context"
	"numoney/internal/bot"
	"numoney/internal/user"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) GetStatistics(ctx context.Context, user user.User, update tgbotapi.Update) error {
	var from time.Time
	var to time.Time
	now := time.Now()

	if update.CallbackQuery == nil {
		return nil
	}

	switch update.CallbackQuery.Data {
	case bot.TodayStatsButton:
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0,
			now.Location())
		to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0,
			now.Location())

	case bot.WeekStatsButton:
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		from = time.Date(now.Year(), now.Month(),
			now.Day()-int(weekday)+1, 0, 0, 0, 0, now.Location())
		to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0,
			now.Location())

	case bot.MonthStatsButton:
		from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0,
			now.Location())
		to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0,
			now.Location())

	case bot.AllStatsButton:
		from = time.Date(2000, 1, 1, 0, 0, 0, 0, now.Location())
		to = time.Now()

	default:
		return nil
	}

	stats, err := h.tService.GetStatsForPeriod(ctx, user.ID, from, to)

	if err != nil {
		return err
	}

	return h.sendEdit(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		bot.Back(),
		bot.GetStatsForPeriodForTg(from, to, stats),
	)
}

func (h *Handler) ChooseStatisticsPeriod(_ context.Context, user user.User, update tgbotapi.Update) error {
	h.states.Set(user.ID, bot.StateStatistics)

	return h.sendEdit(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		bot.ChoseStatPeriod(),
		"Выберите период",
	)
}
