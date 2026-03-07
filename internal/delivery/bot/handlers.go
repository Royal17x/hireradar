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
	logger.Info("Bot responding", "tgID", tgID)
	return c.Send("Привет, это бот с удобным поиском вакансий, рад тебя приветствовать!")
}

func (b *Bot) handleVacancies(c telebot.Context) error {
	ctx := context.Background()
	tgID := strconv.FormatInt(c.Sender().ID, 10)
	user, err := b.userRepo.GetByTelegramID(ctx, tgID)
	if err != nil {
		logger.Error("Failed to get user by telegram ID", "tgID", tgID, "err", err)
		return c.Send("Что-то пошло не так")
	}
	if user == nil {
		logger.Warn("User have not said /start", "tgID", tgID)
		return c.Send("Нажмите /start для начала работы")
	}

	vacancies, err := b.vacancyUcase.GetFiltered(ctx, user.UserID)
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
	logger.Info("Vacancies are successfully showed")
	return c.Send(sb.String(), telebot.ModeMarkdown)
}

func (b *Bot) handleSetFilter(c telebot.Context) error {
	ctx := context.Background()
	tgID := strconv.FormatInt(c.Sender().ID, 10)
	user, err := b.userRepo.GetByTelegramID(ctx, tgID)
	if err != nil {
		logger.Error("Failed to get user by telegram ID", "tgID", tgID, "err", err)
		return c.Send("Не удалось найти вашего пользователя")
	}
	if user == nil {
		logger.Warn("User have not said /start", "tgID", tgID)
		return c.Send("Нажмите /start, для установки фильтра")
	}
	args := c.Args()
	if len(args) == 0 {
		logger.Warn("User haven't set filter", "tgID", tgID)
		return c.Send("Укажи хотя бы ключевое слово. Пример: /setfilter golang москва middle")
	}
	filter := domain.Filter{UserID: user.UserID, CreatedAt: time.Now()}
	filter.Keywords = args[0]
	if len(args) > 1 {
		filter.City = args[1]
	}
	if len(args) > 2 {
		filter.Grade = args[2]
	}

	if err = b.filterRepo.Save(ctx, filter); err != nil {
		logger.Error("Failed to save filter", "tgID", tgID, "err", err)
		return c.Send("Ошибка при сохранении фильтра")
	}
	logger.Info("Filter saved", "tgID", tgID)
	return c.Send("Фильтр успешно записан")
}

func (b *Bot) handleFilter(c telebot.Context) error {
	ctx := context.Background()
	tgID := strconv.FormatInt(c.Sender().ID, 10)
	user, err := b.userRepo.GetByTelegramID(ctx, tgID)
	if err != nil {
		logger.Error("Failed to get user by telegram ID", "tgID", tgID, "err", err)
		return c.Send("Не удалось найти вашего пользователя")
	}
	if user == nil {
		logger.Warn("User have not said /start", "tgID", tgID)
		return c.Send("Нажмите /start, для работы бота")
	}
	filters, err := b.filterRepo.GetByUserID(ctx, user.UserID)
	if err != nil {
		logger.Error("Failed to get filters", "tgID", tgID, "err", err)
		return c.Send("Не удалось вывести фильтры")
	}
	if len(filters) == 0 {
		return c.Send("У вас нет активных фильтров")
	}

	var sb strings.Builder
	for i, filter := range filters {
		sb.WriteString(fmt.Sprintf("*Фильтр №%d* (ID: %d)\nКлючевые слова: %s\nГород: %s\nГрейд: %s\nСоздан: %s\n\n",
			i+1, filter.ID, filter.Keywords, filter.City, filter.Grade, filter.CreatedAt.Format("02.01.2006")))
	}
	logger.Info("Filters are successfully showed")
	return c.Send(sb.String(), telebot.ModeMarkdown)
}

func (b *Bot) handleDeleteFilter(c telebot.Context) error {
	ctx := context.Background()
	args := c.Args()
	if len(c.Args()) == 0 {
		logger.Warn("Wrong parametres", "tgID", c.Sender().ID)
		return c.Send("Укажи ID фильтра. Пример: /deletefilter 3")
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Error("Failed to parse id", "args", args, "err", err)
		return c.Send("Неправильно заданный параметр")
	}
	if err = b.filterRepo.Delete(ctx, id); err != nil {
		logger.Error("Failed to delete filter", "id", id, "err", err)
		return c.Send("Произошла ошибка удаления фильтра")
	}
	logger.Info("Filter successfully deleted", "id", id)
	return c.Send("Фильтр успешно удален")
}
