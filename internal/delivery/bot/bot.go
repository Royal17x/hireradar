package bot

import (
	"github.com/Royal17x/hireradar/internal/domain"
	"github.com/Royal17x/hireradar/internal/usecase"
	"gopkg.in/telebot.v3"
	"time"
)

type Bot struct {
	bot          *telebot.Bot
	vacancyUcase *usecase.VacancyUsecase
	userRepo     domain.UserRepository
	filterRepo   domain.FilterRepository
}

func NewBot(token string, vacancyUcase *usecase.VacancyUsecase, userRepo domain.UserRepository, filterRepo domain.FilterRepository) (*Bot, error) {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}
	bot := &Bot{
		bot:          b,
		vacancyUcase: vacancyUcase,
		userRepo:     userRepo,
		filterRepo:   filterRepo,
	}

	b.Handle("/start", bot.handleStart)
	b.Handle("/vacancies", bot.handleVacancies)
	return bot, nil
}

func (b *Bot) Start() {
	b.bot.Start()
}
