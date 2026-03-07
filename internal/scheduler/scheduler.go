package scheduler

import (
	"context"
	"github.com/Royal17x/hireradar/internal/usecase"
	logger "github.com/charmbracelet/log"
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
	logger.Infof("Starting scheduler with interval %v", s.interval)
	err := s.usecase.FetchAndStore(ctx, s.query)
	if err != nil {
		logger.Error("Fetch and store vacancies error", "err", err)
	}
	for {
		select {
		case <-ticker.C:
			err = s.usecase.FetchAndStore(ctx, s.query)
			if err != nil {
				logger.Error("Fetch and store vacancies error", "err", err)
			}
		case <-ctx.Done():
			logger.Info("Shutting down scheduler...")
			return
		}

	}

}
