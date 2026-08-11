package transaction

type Transaction struct {
	ID         int
	Amount     int
	CategoryID int
	UserID     int64
	Comment    string
}
