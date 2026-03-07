package usecase

import (
	"context"
	"errors"
	"github.com/Royal17x/hireradar/internal/domain"
	logger "github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgconn"
)

type VacancyUsecase struct {
	vacancyRepo domain.VacancyRepository
	cacheRepo   domain.CacheRepository
	fetcher     domain.VacancyFetcher
}

func NewVacancyUsecase(vacancyRepo domain.VacancyRepository, cacheRepo domain.CacheRepository, fetcher domain.VacancyFetcher) *VacancyUsecase {
	return &VacancyUsecase{
		vacancyRepo: vacancyRepo,
		cacheRepo:   cacheRepo,
		fetcher:     fetcher,
	}
}

func (u *VacancyUsecase) FetchAndStore(ctx context.Context, query string) error {
	vacancies, err := u.fetcher.FetchVacancies(ctx, query)
	if err != nil {
		logger.Error("Fetch vacancies error", "err", err)
		return err
	}
	for _, vacancy := range vacancies {
		seen, err := u.cacheRepo.IsSeen(ctx, vacancy.HhID)
		if err != nil {
			logger.Error("Error checking vacancy existence", "err", err)
			return err
		}
		if seen {
			logger.Debug("Vacancy already seen", "hh_id", vacancy.HhID)
			continue
		}
		v := vacancy
		if err = u.vacancyRepo.Save(ctx, &v); err != nil {
			if isUniqueViolation(err) {
				logger.Info("Redis have not found hh_id after restart")
				continue
			}
			logger.Error("Error storing vacancy", "err", err)
			return err
		}
		err = u.cacheRepo.SetSeen(ctx, vacancy.HhID)
		if err != nil {
			logger.Error("Failed to set vacancy to cache", "hh_id", vacancy.HhID)
			return err
		}
	}
	return nil
}

func (u *VacancyUsecase) GetAll(ctx context.Context) ([]domain.Vacancy, error) {
	return u.vacancyRepo.GetAll(ctx)
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
