package domain

import "context"

type VacancyRepository interface {
	Save(ctx context.Context, vacancy *Vacancy) error
	Exists(ctx context.Context, hhID string) (bool, error)
	GetAll(ctx context.Context) ([]Vacancy, error)
}

type UserRepository interface {
	Save(ctx context.Context, user User) error
	GetByTelegramID(ctx context.Context, tgID string) (*User, error)
}
