package scheduler

import (
	"context"
	"github.com/Royal17x/hireradar/internal/usecase"
	"log"
	"time"
)

type Scheduler struct {
	usecase  *usecase.VacancyUsecase
	interval time.Duration
	query    string
}

func NewScheduler(usecase *usecase.VacancyUsecase, interval time.Duration, query string) *Scheduler {
	return &Scheduler{
		usecase:  usecase,
		interval: interval,
		query:    query,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	err := s.usecase.FetchAndStore(ctx, s.query)
	if err != nil {
		log.Printf("ошибка сбора и хранения вакансий: %v", err)
	}
	for {
		select {
		case <-ticker.C:
			err = s.usecase.FetchAndStore(ctx, s.query)
			if err != nil {
				log.Printf("ошибка сбора и хранения вакансий: %v", err)
			}
		case <-ctx.Done():
			return
		}

	}

}
