package bot

type State struct {
	state            map[int64]string
	selectedCategory map[int64]int
}

func NewState() *State {
    return &State{
        state:            make(map[int64]string),
        selectedCategory: make(map[int64]int),
    }
}

const (
	StateNone                 = ""
	StateAddCategoryName      = "add_category_name"
	StateAddCategoryType      = "add_category_type"
	StateAddTransaction       = "add_transaction"
	StateAddTransactionAmount = "add_transaction_amount"
	StateStaistics            = "statistics"
)
