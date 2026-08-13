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
	uService  user.Service
	state     *State
}

func NewRouter(uService user.Service, state *State) *Router {
	return &Router{
		callbacks: make(map[string]HandlerFunc),
		states:    make(map[string]HandlerFunc),
		uService:  uService,
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

	user, err := r.uService.EnsureRegistered(
		ctx,
		update.CallbackQuery.From.ID,
	)

	if update.CallbackQuery != nil {
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
