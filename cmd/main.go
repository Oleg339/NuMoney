package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"numoney/internal/category"
	"numoney/internal/transaction"
	"numoney/internal/user"
	"os"
	"strconv"
	"strings"

	"time"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const (
	StateNone                 = ""
	StateAddCategoryName      = "add_category_name"
	StateAddCategoryType      = "add_category_type"
	StateAddTransaction       = "add_transaction"
	StateAddTransactionAmount = "add_transaction_amount"
	StateStaistics            = "statistics"
)

func getBaseMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
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

func main() {
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")

	db := getConnection()
	defer db.Close()

	categoryRepository := category.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)
	userRepository := user.NewRepository(db)

	var userStates = map[int64]string{}
	var userSelectedCategory = map[int64]int{}

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug, _ = strconv.ParseBool(os.Getenv("BOT_DEBUG"))

	log.Printf("Authorized on account %s", bot.Self.UserName)

	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60

	updates := bot.GetUpdatesChan(update)

	var u user.User

	for update := range updates {
		var text string
		var keyboard tgbotapi.InlineKeyboardMarkup
		var msg tgbotapi.EditMessageTextConfig

		if update.Message != nil {
			text = "Выберите действие:"
			keyboard = getBaseMenuKeyboard()

			u, err = userRepository.GetByTelegramId(update.Message.From.ID)

			if err != nil {
				log.Fatal(err)
			}

			switch userStates[int64(u.ID)] {
			case "add_expense_category", "add_income_category":
				catType := "income"

				if userStates[int64(u.ID)] == "add_expense_category" {
					catType = "expense"
				}

				categoryRepository.Save(&category.Category{
					Name:   update.Message.Text,
					UserId: u.ID,
					Type:   catType,
				})
			case StateAddTransactionAmount:
				text = strings.Replace(update.Message.Text, ",", ".", 1)
				amount, err := strconv.ParseFloat(text, 64)

				if err != nil {
					text = "Введите корректное число:"
					userStates[int64(u.ID)] = StateAddTransactionAmount
					keyboard = tgbotapi.InlineKeyboardMarkup{}

				} else {
					transactionRepository.Save(&transaction.Transaction{
						Amount:     int(amount * 100),
						CategoryId: userSelectedCategory[int64(u.ID)],
						UserId:     u.ID,
					})

					userSelectedCategory[int64(u.ID)] = 0
					userStates[int64(u.ID)] = StateNone
				}

			}

			msgNew := tgbotapi.NewMessage(update.Message.Chat.ID, text)

			if len(keyboard.InlineKeyboard) > 0 {
				msgNew.ReplyMarkup = keyboard
			}

			bot.Send(msgNew)
			continue
		}

		if update.CallbackQuery != nil {

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			u, err = userRepository.GetByTelegramId(update.CallbackQuery.From.ID)

			if err != nil {
				log.Fatal(err)
			}

			bot.Request(callback)

			if userStates[int64(u.ID)] == StateStaistics {
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

				stats, err := transactionRepository.GetStatsForPeriodFromDb(u.ID, from, to)

				text = getStatsForPeriodForTg(from, to, stats)

				if err != nil {
					log.Fatal(err)
				}

				keyboard = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back"),
					),
				)
			}

			switch update.CallbackQuery.Data {
			case "add_category":
				text = "Выберите тип"
				keyboard = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back"),
					),
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🔴 Расход", "add_expense_category"),
						tgbotapi.NewInlineKeyboardButtonData("🟢 Доход", "add_income_category"),
					),
				)
			case "Statistics":
				text = "Выберите период"
				keyboard = tgbotapi.NewInlineKeyboardMarkup(
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

				userStates[int64(u.ID)] = StateStaistics
			case "add_expense_category", "add_income_category":
				userStates[int64(u.ID)] = update.CallbackQuery.Data

				text = "Введите название категории"

				keyboard = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back"),
					),
				)
			case "back":
				text = "Выберите действие:"
				keyboard = getBaseMenuKeyboard()

			case "Categories":
				categories, err := categoryRepository.GetByUserId(u.ID)

				if err != nil {
					log.Fatal(err)
				}

				if len(categories) == 0 {
					text = "У вас нет категорий."
					keyboard = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("➕ Добавить", "add_category"),
							tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back"),
						),
					)
				} else {
					text = "Категории"
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

					for _, category := range categories {
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
						log.Fatal(err)
					}
					resp, err := bot.MakeRequest("editMessageReplyMarkup", tgbotapi.Params{
						"chat_id":      strconv.FormatInt(int64(u.TelegramId), 10),
						"message_id":   fmt.Sprint(update.CallbackQuery.Message.MessageID),
						"text":         text,
						"reply_markup": string(markup),
					})

					fmt.Print(resp)
					continue
				}

			case "Expense", "Income":
				text := "Выберите категорию"

				userStates[int64(u.ID)] = StateAddTransaction

				categories, err := categoryRepository.GetByUserIdAndType(u.ID, strings.ToLower(update.CallbackQuery.Data))

				if err != nil {
					log.Fatal(err)
				}

				markup := getCategoriesMarkup(categories)

				resp, err := bot.MakeRequest("editMessageReplyMarkup", tgbotapi.Params{
					"chat_id":      strconv.FormatInt(int64(u.TelegramId), 10),
					"message_id":   fmt.Sprint(update.CallbackQuery.Message.MessageID),
					"text":         text,
					"reply_markup": string(markup),
				})

				if err != nil {
					log.Fatal(err)
				}

				fmt.Print(resp)
				continue

			}

			if userStates[int64(u.ID)] == StateAddTransaction {
				categories, err := categoryRepository.GetByUserId(u.ID)

				if err != nil {
					log.Fatal(err)
				}

				for _, cat := range categories {
					if fmt.Sprint(cat.ID) == update.CallbackQuery.Data {
						userSelectedCategory[int64(u.ID)] = cat.ID
						userStates[int64(u.ID)] = StateAddTransactionAmount
						text = "Введите сумму:"
						msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, text)

						bot.Send(msg)
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

		bot.Send(msg)
	}
}

func getConnection() *sql.DB {
	cStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	db, err := sql.Open(
		"mysql",
		cStr,
	)

	log.Print(cStr)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

type Transaction struct {
	ID         int
	Amount     int
	CategoryId int
	UserId     int
	Comment    string
}

func getCategoriesMarkup(categories []category.Category) []byte {
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

	for _, category := range categories {
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
		log.Fatal(err)
	}

	return markup
}

func getStatsForPeriodForTg(from time.Time, to time.Time, stats []transaction.Stat) string {
	expensesSum := 0
	incomeSum := 0
	total := 0

	for _, stat := range stats {
		if stat.Type == "expense" {
			expensesSum += stat.Amount
			total -= stat.Amount
			continue
		}

		incomeSum += stat.Amount
		total += stat.Amount
	}

	text := fmt.Sprintf(
		"Статистика за период %s - %s\n\n🔴 Расходы: <code>%.2f</code>",
		from.Format("02.01.2006"),
		to.Format("02.01.2006"),
		float64(expensesSum)/100,
	)

	for _, stat := range stats {
		if stat.Type == "expense" {
			text = fmt.Sprint(text, fmt.Sprintf("\n\t\t\t* %s - <code>%.2f</code>", stat.Name, float64(stat.Amount)/100))
		}
	}

	text = fmt.Sprint(text, fmt.Sprintf("\n\n🟢 Доходы: <code>%.2f</code>", float64(incomeSum)/100))

	for _, stat := range stats {
		if stat.Type == "income" {
			text = fmt.Sprint(text, fmt.Sprintf("\n\t\t\t* %s - <code>%.2f</code>", stat.Name, float64(stat.Amount)/100))
		}
	}

	return fmt.Sprint(text, fmt.Sprintf("\n\nБаланс: <code>%.2f</code>", float64(total)/100))
}
