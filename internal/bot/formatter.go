package bot

import (
	"fmt"
	"numoney/internal/transaction"
	"strings"
	"time"
)

func GetStatsForPeriodForTg(from time.Time, to time.Time, stats []transaction.Stat) string {
	var expensesSum, incomeSum, total int
	var expenses, incomes strings.Builder

	for _, stat := range stats {
		line := fmt.Sprintf("\n\t\t\t* %s - <code>%.2f</code>", stat.Name, float64(stat.Amount)/100)
		if stat.Type == "expense" {
			expensesSum += stat.Amount
			total -= stat.Amount
			expenses.WriteString(line)
		} else {
			incomeSum += stat.Amount
			total += stat.Amount
			incomes.WriteString(line)
		}
	}

	return fmt.Sprintf(
		"Статистика за период %s - %s\n\n🔴 Расходы: <code>%.2f</code>%s\n\n🟢 Доходы: <code>%.2f</code>%s\n\nБаланс: <code>%.2f</code>",
		from.Format("02.01.2006"),
		to.Format("02.01.2006"),
		float64(expensesSum)/100,
		expenses.String(),
		float64(incomeSum)/100,
		incomes.String(),
		float64(total)/100,
	)
}
