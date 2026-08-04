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

func (s *State) Get(userId int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state[userId]
}

func (s *State) Set(userId int64, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state[userId] = state
}

func (s *State) GetCategory(userId int64) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.selectedCategory[userId]
}

func (s *State) SetCategory(userId int64, category int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.selectedCategory[userId] = category
}

const (
	StateNone                 = ""
	StateAddCategoryName      = "add_category_name"
	StateAddCategoryType      = "add_category_type"
	StateAddTransaction       = "add_transaction"
	StateAddTransactionAmount = "add_transaction_amount"
	StateStaistics            = "statistics"
)
