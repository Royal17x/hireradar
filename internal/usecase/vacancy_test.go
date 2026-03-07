package usecase

import (
	"context"
	"errors"
	"github.com/Royal17x/hireradar/internal/domain"
	"github.com/Royal17x/hireradar/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFetch(t *testing.T) {
	tests := []struct {
		name          string
		hhID          string
		wantSaveCalls int
		isSeen        bool
		wantErr       bool
		fetcherErr    error
		saveErr       error
		isSeenErr     error
	}{
		{
			name:       "ошибка fetcher",
			fetcherErr: errors.New("ошибка клиента"),
			wantErr:    true,
		},
		{
			name:          "ошибка сохранения в БД",
			saveErr:       errors.New("ошибка БД"),
			wantSaveCalls: 1,
			wantErr:       true,
		},
		{
			name:    "вакансия уже была сохранена",
			hhID:    "123",
			isSeen:  true,
			wantErr: false,
		},
		{
			name:          "новая вакансия пришла из парсинга",
			hhID:          "123",
			wantSaveCalls: 1,
			isSeen:        false,
			wantErr:       false,
		},
		{
			name:      "ошибка проверки вакансии из существующих",
			hhID:      "123",
			isSeenErr: errors.New("ошибка БД"),
			wantErr:   true,
		},
	}

	// arrange
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vacancyRepo := mocks.NewVacancyRepository(t)
			cacheRepo := mocks.NewCacheRepository(t)
			filterRepo := mocks.NewFilterRepository(t)
			fetcher := mocks.NewVacancyFetcher(t)

			fetcher.On("FetchVacancies", mock.Anything, "golang").
				Return([]domain.Vacancy{{HhID: tt.hhID}}, tt.fetcherErr)
			if tt.fetcherErr == nil {
				cacheRepo.On("IsSeen", mock.Anything, tt.hhID).
					Return(tt.isSeen, tt.isSeenErr)
			}
			if tt.wantSaveCalls > 0 && !tt.isSeen && tt.isSeenErr == nil {
				vacancyRepo.On("Save", mock.Anything, mock.Anything).
					Return(tt.saveErr)
				if tt.saveErr == nil {
					cacheRepo.On("SetSeen", mock.Anything, tt.hhID).
						Return(tt.saveErr)
				}

			}

			// act
			ucase := NewVacancyUsecase(vacancyRepo, cacheRepo, filterRepo, fetcher)
			err := ucase.FetchAndStore(context.Background(), "golang")

			// assert
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}

}
