package bot

import "sync"

type State struct {
	mu               sync.RWMutex
	state            map[int64]string
	selectedCategory map[int64]int
}

func NewState() *State {
	return &State{
		state:            make(map[int64]string),
		selectedCategory: make(map[int64]int),
	}
}

func (s *State) Get(userID int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state[userID]
}

func (s *State) Set(userID int64, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state[userID] = state
}

func (s *State) GetCategory(userID int64) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.selectedCategory[userID]
}

func (s *State) SetCategory(userID int64, category int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.selectedCategory[userID] = category
}

const (
	StateNone                     = ""
	StateAddTransactionAmount     = "add_transaction_amount"
	StateStatistics                = "statistics"
	StateAddExpenseCategory       = "add_expense_category"
	StateAddIncomeCategory        = "add_income_category"
	BackButton                    = "back"
	StateChoseTransactionCategory = "chose_transaction_category"
)
