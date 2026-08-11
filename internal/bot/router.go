package bot

import (
	"context"
	"time"

	"numoney/internal/user"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerFunc func(ctx context.Context, user user.User, update tgbotapi.Update) error

type Router struct {
	callbacks map[string]HandlerFunc
	states    map[string]HandlerFunc
	uRepo     user.Repository
	state     *State
}

func NewRouter(uRepo user.Repository, state *State) *Router {
	return &Router{
		callbacks: make(map[string]HandlerFunc),
		states:    make(map[string]HandlerFunc),
		uRepo:     uRepo,
		state:     state,
	}
}

func (r *Router) OnCallback(data string, f HandlerFunc) {
	r.callbacks[data] = f
}

func (r *Router) OnState(data string, f HandlerFunc) {
	r.states[data] = f
}

func (r *Router) Handle(update tgbotapi.Update) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if update.CallbackQuery != nil {
		user, err := r.uRepo.GetByTelegramID(ctx, update.CallbackQuery.From.ID)

		if err != nil {
			return err
		}

		fn, ok := r.callbacks[update.CallbackQuery.Data]

		if ok {
			return fn(ctx, user, update)
		}

		state := r.state.Get(user.ID)

		fn, ok = r.states[state]

		if ok {
			return fn(ctx, user, update)
		}
	}

	if update.Message != nil {
		user, err := r.uRepo.GetByTelegramID(ctx, update.Message.From.ID)

		if err != nil {
			return err
		}

		state := r.state.Get(user.ID)
		
		fn, ok := r.states[state]
		if ok {
			return fn(ctx, user, update)
		}
	}

	return nil
}
