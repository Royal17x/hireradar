package usecase

import (
	"context"
	"github.com/Royal17x/hireradar/internal/domain"
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
		return err
	}
	for _, vacancy := range vacancies {
		seen, err := u.cacheRepo.IsSeen(ctx, vacancy.HhID)
		if err != nil {
			return err
		}
		if seen {
			continue
		}
		v := vacancy
		if err = u.vacancyRepo.Save(ctx, &v); err != nil {
			return err
		}
		err = u.cacheRepo.SetSeen(ctx, vacancy.HhID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *VacancyUsecase) GetAll(ctx context.Context) ([]domain.Vacancy, error) {
	return u.vacancyRepo.GetAll(ctx)
}

func isUniqueViolation(err error) bool {
	//TODO: skip constraint error (after redis restart)
	return true
}
