package main

import (
	"database/sql"
	"fmt"
	"log"
	"numoney/internal/bot"
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

	categoryRepository := category.NewMySQLRepository(db)
	transactionRepository := transaction.NewMySQLRepository(db)
	userRepository := user.NewMySQLRepository(db)
	states := bot.NewState()

	b, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Panic(err)
	}

	b.Debug, _ = strconv.ParseBool(os.Getenv("BOT_DEBUG"))

	log.Printf("Authorized on account %s", b.Self.UserName)

	update := tgbotapi.NewUpdate(0)

	update.Timeout = 60

	updates := b.GetUpdatesChan(update)

	handler := bot.NewHandler(b, userRepository, categoryRepository, transactionRepository, states)

	for update := range updates {
		handler.Handle(update)
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

	return db
}
