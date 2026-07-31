package bot

import (
	"fmt"
	"log"
	"numoney/internal/category"
	"numoney/internal/transaction"
	"numoney/internal/user"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot    *tgbotapi.BotAPI
	uRepo  user.Repository
	cRepo  category.Repository
	tRepo  transaction.Repository
	states *State
}

func NewHandler(bot *tgbotapi.BotAPI, uRepo user.Repository, cRepo category.Repository, tRepo transaction.Repository, states *State) Handler {
	return Handler{bot, uRepo, cRepo, tRepo, states}
}

func (h *Handler) Handle(update tgbotapi.Update) error {
	var states = h.states
	text := ""
	var keyboard tgbotapi.InlineKeyboardMarkup
	var msg tgbotapi.EditMessageTextConfig

	fmt.Println(states)

	if update.Message != nil {
		text = "Выберите действие:"

		keyboard = BaseMenu()

		u, err := h.uRepo.GetByTelegramId(update.Message.From.ID)

		if err != nil {
			return err
		}

		userId := int64(u.ID)

		switch states.state[userId] {

		case "add_expense_category", "add_income_category":
			catType := "income"

			if states.state[userId] == "add_expense_category" {
				catType = "expense"
			}

			h.cRepo.Save(&category.Category{
				Name:   update.Message.Text,
				UserId: u.ID,
				Type:   catType,
			})
		case StateAddTransactionAmount:
			text = strings.Replace(update.Message.Text, ",", ".", 1)
			amount, err := strconv.ParseFloat(text, 64)

			if err != nil {
				text = "Введите корректное число:"

				states.state[userId] = StateAddTransactionAmount
				keyboard = tgbotapi.InlineKeyboardMarkup{}

			} else {
				h.tRepo.Save(&transaction.Transaction{
					Amount:     int(amount * 100),
					CategoryId: states.selectedCategory[userId],
					UserId:     u.ID,
				})

				states.selectedCategory[userId] = 0
				states.state[userId] = StateNone
			}

		}

		msgNew := tgbotapi.NewMessage(update.Message.Chat.ID, text)

		if len(keyboard.InlineKeyboard) > 0 {
			msgNew.ReplyMarkup = keyboard
		}

		h.bot.Send(msgNew)
	}

	if update.CallbackQuery != nil {

		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		u, err := h.uRepo.GetByTelegramId(update.CallbackQuery.From.ID)
		userId := int64(u.ID)

		if err != nil {
			log.Fatal(err)
		}

		h.bot.Request(callback)

		if states.state[userId] == StateStaistics {
			var from time.Time
			var to time.Time
			now := time.Now()

			if update.CallbackQuery.Data == "today_stats" {
				from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0,
					now.Location())
				to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0,
					now.Location())
			}

			if update.CallbackQuery.Data == "week_stats" {
				from = time.Date(now.Year(), now.Month(),
					now.Day()-int(now.Weekday())+1, 0, 0, 0, 0, now.Location())
				to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0,
					now.Location())
			}

			if update.CallbackQuery.Data == "month_stats" {
				from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0,
					now.Location())
				to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0,
					now.Location())
			}

			if update.CallbackQuery.Data == "all_stats" {
				from = time.Date(now.Year()-1000, now.Month(),
					0, 0, 0, 0, 0, now.Location())
				to = time.Date(now.Year()+1000, now.Month(), now.Day(), 23, 59, 59, 0,
					now.Location())
			}

			stats, err := h.tRepo.GetStatsForPeriodFromDb(u.ID, from, to)

			text = getStatsForPeriodForTg(from, to, stats)

			if err != nil {
				return err
			}

			keyboard = Back()
		}

		switch update.CallbackQuery.Data {
		case "add_category":
			text = "Выберите тип"
			keyboard = AddCategory()

		case "Statistics":
			text = "Выберите период"
			keyboard = ChoseStatPeriod()
			states.state[userId] = StateStaistics

		case "add_expense_category", "add_income_category":
			states.state[userId] = update.CallbackQuery.Data
			text = "Введите название категории"
			keyboard = Back()

		case "back":
			text = "Выберите действие:"
			keyboard = BaseMenu()
			states.state[userId] = StateNone

		case "Categories":
			categories, err := h.cRepo.GetByUserId(u.ID)

			if err != nil {
				log.Fatal(err)
			}

			if len(categories) == 0 {
				text = "У вас нет категорий."
				keyboard = NoCategories()
			} else {
				text = "Категории"

				markup, err := CategoriesJSON(categories)

				if err != nil {
					log.Fatal(err)
				}
				resp, err := h.bot.MakeRequest("editMessageReplyMarkup", tgbotapi.Params{
					"chat_id":      strconv.FormatInt(int64(u.TelegramId), 10),
					"message_id":   fmt.Sprint(update.CallbackQuery.Message.MessageID),
					"text":         text,
					"reply_markup": string(markup),
				})

				fmt.Print(resp)

				return nil
			}

		case "Expense", "Income":
			text = "Выберите категорию"

			states.state[userId] = StateAddTransaction

			categories, err := h.cRepo.GetByUserIdAndType(u.ID, strings.ToLower(update.CallbackQuery.Data))

			if err != nil {
				log.Fatal(err)
			}

			markup, err := ChoseCategoryJSON(categories)

			if err != nil {
				return err
			}

			resp, err := h.bot.MakeRequest("editMessageReplyMarkup", tgbotapi.Params{
				"chat_id":      strconv.FormatInt(int64(u.TelegramId), 10),
				"message_id":   fmt.Sprint(update.CallbackQuery.Message.MessageID),
				"text":         text,
				"reply_markup": string(markup),
			})

			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(resp)

			return nil
		}

		if states.state[userId] == StateAddTransaction {
			categories, err := h.cRepo.GetByUserId(u.ID)

			if err != nil {
				log.Fatal(err)
			}

			for _, cat := range categories {
				if fmt.Sprint(cat.ID) == update.CallbackQuery.Data {
					states.selectedCategory[userId] = cat.ID
					states.state[userId] = StateAddTransactionAmount
					text = "Введите сумму:"
					msg = tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, text)

					h.bot.Send(msg)
					break
				}
			}
		}

		msg = tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, text)
		msg.ParseMode = "HTML"
	}

	if len(keyboard.InlineKeyboard) > 0 {
		msg.ReplyMarkup = &keyboard
	}

	h.bot.Send(msg)

	return nil
}
