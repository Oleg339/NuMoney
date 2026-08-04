package transaction

type Transaction struct {
	ID         int
	Amount     int
	CategoryId int
	UserId     int
	Comment    string
}
