package main

import (
	"database/sql"
	"fmt"
	"log"
	"numoney/internal/bot"
	"numoney/internal/bot/handlers"
	"numoney/internal/category"
	"numoney/internal/transaction"
	"numoney/internal/user"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")

	db := getConnection()
	defer db.Close()

	states := bot.NewState()

	b, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Panic(err)
	}

	b.Debug, _ = strconv.ParseBool(os.Getenv("BOT_DEBUG"))

	cfg := tgbotapi.NewUpdate(0)

	cfg.Timeout = 60

	updates := b.GetUpdatesChan(cfg)

	handler := handlers.NewHandler(
		b,
		category.NewService(
			category.NewMySQLRepository(db),
		),
		transaction.NewService(
			transaction.NewMySQLRepository(db),
			category.NewMySQLRepository(db),
		),
		states,
	)

	router := bot.NewRouter(
		user.NewService(
			user.NewMySQLRepository(db),
		),
		states,
	)

	router.OnState(bot.StateNone, handler.InitMessage)
	router.OnCallback(bot.BackButton, handler.InitMessage)

	router.OnCallback(bot.ChoseStatisticsButton, handler.ChooseStatisticsPeriod)
	router.OnState(bot.StateStatistics, handler.GetStatistics)

	router.OnCallback(bot.ChoseCategoryButton, handler.ChoseCategory)
	router.OnCallback(bot.AddCategoryButton, handler.AddCategoryButton)
	router.OnState(bot.ChoseCategoryToCreateButton, handler.ChoseCategoryToCreateButton)
	router.OnState(bot.StateAddExpenseCategory, handler.AddCategory)
	router.OnState(bot.StateAddIncomeCategory, handler.AddCategory)

	router.OnCallback(bot.AddExpenseButton, handler.AddTransactionButton)
	router.OnCallback(bot.AddIncomeButton, handler.AddTransactionButton)
	router.OnState(bot.StateChoseTransactionCategory, handler.StateChoseTransactionCategory)
	router.OnState(bot.StateAddTransactionAmount, handler.AddTransactionAmount)

	for update := range updates {
		go func(u tgbotapi.Update) {
			err = router.Handle(u)

			if err != nil {
				log.Println("handle error:", err)
			}
		}(update)
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

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	fmt.Println("connections: ", db.Stats().OpenConnections)

	return db
}
