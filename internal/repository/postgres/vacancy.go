package postgres

import (
	"context"
	"fmt"
	"github.com/Royal17x/hireradar/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type VacancyRepository struct {
	db *pgxpool.Pool
}

func NewVacancyRepository(db *pgxpool.Pool) *VacancyRepository {
	return &VacancyRepository{db: db}
}

func (v *VacancyRepository) Save(ctx context.Context, vacancy *domain.Vacancy) error {

	query := `
INSERT INTO vacancies 
    (hh_id, title, city, company, url, salary_from, salary_to, published_at, created_at) 
VALUES 
    ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	_, err := v.db.Exec(ctx, query, vacancy.HhID, vacancy.Title, vacancy.City, vacancy.Company, vacancy.URL, vacancy.SalaryFrom, vacancy.SalaryTo, vacancy.PublishedAt, vacancy.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (v *VacancyRepository) Exists(ctx context.Context, hhID string) (bool, error) {
	query := `SELECT COUNT(*)
FROM vacancies
WHERE hh_id = $1;`
	var count int
	res := v.db.QueryRow(ctx, query, hhID).Scan(&count)
	if res != nil {
		return false, nil
	}
	if count <= 0 {
		return false, nil
	}
	return true, nil
}

func (v *VacancyRepository) GetAll(ctx context.Context) ([]domain.Vacancy, error) {
	query := `SELECT  hh_id, title, city, company, url, salary_from, salary_to, published_at, created_at
FROM vacancies;`
	rows, err := v.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var vacancies []domain.Vacancy
	for rows.Next() {
		vacancy := domain.Vacancy{}
		rows.Scan(&vacancy.HhID, &vacancy.Title, &vacancy.City, &vacancy.Company, &vacancy.URL, &vacancy.SalaryFrom, &vacancy.SalaryTo, &vacancy.PublishedAt, &vacancy.CreatedAt)
		vacancies = append(vacancies, vacancy)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return vacancies, nil
}

func (v *VacancyRepository) GetFiltered(ctx context.Context, keywords, city, grade string) ([]domain.Vacancy, error) {
	conditions := []string{"1=1"}
	args := []any{}
	argIdx := 1

	if keywords != "" {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", argIdx))
		args = append(args, "%"+keywords+"%")
		argIdx++
	}
	if city != "" {
		conditions = append(conditions, fmt.Sprintf("city ILIKE $%d", argIdx))
		args = append(args, "%"+city+"%")
		argIdx++
	}
	if grade != "" {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", argIdx))
		args = append(args, "%"+grade+"%")
		argIdx++
	}
	query := `SELECT  hh_id, title, city, company, url, salary_from, salary_to, published_at, created_at
FROM vacancies
WHERE ` + strings.Join(conditions, " AND ")

	var vacancies []domain.Vacancy
	rows, err := v.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		vacancy := domain.Vacancy{}
		rows.Scan(&vacancy.HhID, &vacancy.Title, &vacancy.City, &vacancy.Company, &vacancy.URL, &vacancy.SalaryFrom, &vacancy.SalaryTo, &vacancy.PublishedAt, &vacancy.CreatedAt)
		vacancies = append(vacancies, vacancy)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return vacancies, nil
}
