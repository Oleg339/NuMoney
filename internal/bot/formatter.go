package bot

import (
	"fmt"
	"numoney/internal/transaction"
	"time"
)

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
