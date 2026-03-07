package bot

import (
	"context"
	"fmt"
	"github.com/Royal17x/hireradar/internal/domain"
	"gopkg.in/telebot.v3"
	"strconv"
	"strings"
	"time"
)

func (b *Bot) handleStart(c telebot.Context) error {
	ctx := context.Background()
	tgID := strconv.FormatInt(c.Sender().ID, 10)
	tgUs := c.Sender().Username
	user, err := b.userRepo.GetByTelegramID(ctx, tgID)
	if err != nil {
		return c.Send("произошла ошибка, попробуйте позже")
	}
	if user == nil {
		user = &domain.User{
			TgID:      tgID,
			Username:  tgUs,
			IsActive:  true,
			CreatedAt: time.Now(),
		}
		if err = b.userRepo.Save(ctx, *user); err != nil {
			return c.Send("произошла ошибка пре регистрации")
		}
	}
	return c.Send("привет, это бот с удобным поиском вакансий, рад тебя приветствовать!")
}

func (b *Bot) handleVacancies(c telebot.Context) error {
	vacancies, err := b.vacancyUcase.GetAll(context.Background())
	if err != nil {
		return c.Send("ошибка при получении вакансий")
	}
	if len(vacancies) == 0 {
		return c.Send("вакансий пока нет")
	}

	var sb strings.Builder
	for _, v := range vacancies[:min(len(vacancies), 10)] {
		sb.WriteString(fmt.Sprintf("*%s* - %s\n%s\n\n", v.Title, v.Company, v.URL))
	}
	return c.Send(sb.String(), telebot.ModeMarkdown)
}
