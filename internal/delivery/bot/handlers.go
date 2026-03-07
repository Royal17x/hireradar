package bot

import (
	"context"
	"fmt"
	"github.com/Royal17x/hireradar/internal/domain"
	logger "github.com/charmbracelet/log"
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
		logger.Error("Failed to get user by telegram ID", "tgID", tgID, "err", err)
		return c.Send("Произошла ошибка, попробуйте позже")
	}
	if user == nil {
		user = &domain.User{
			TgID:      tgID,
			Username:  tgUs,
			IsActive:  true,
			CreatedAt: time.Now(),
		}
		if err = b.userRepo.Save(ctx, *user); err != nil {
			logger.Error("Failed to save user", "tgID", tgID, "err", err)
			return c.Send("Произошла ошибка пре регистрации")
		}
		logger.Info("User created", "tgID", tgID)
	}
	return c.Send("Привет, это бот с удобным поиском вакансий, рад тебя приветствовать!")
}

func (b *Bot) handleVacancies(c telebot.Context) error {
	vacancies, err := b.vacancyUcase.GetAll(context.Background())
	if err != nil {
		logger.Error("Failed to get all vacancies", "err", err)
		return c.Send("Ошибка при получении вакансий")
	}
	if len(vacancies) == 0 {
		logger.Warn("No vacancies found")
		return c.Send("Вакансий пока нет")
	}

	var sb strings.Builder
	for _, v := range vacancies[:min(len(vacancies), 10)] {
		sb.WriteString(fmt.Sprintf("*%s* - %s\n%s\n\n", v.Title, v.Company, v.URL))
	}
	return c.Send(sb.String(), telebot.ModeMarkdown)
}
