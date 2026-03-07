package domain

import "time"

type Vacancy struct {
	VacancyID   int
	HhID        string
	Title       string
	City        string
	Company     string
	URL         string
	SalaryFrom  *int
	SalaryTo    *int
	PublishedAt time.Time
	CreatedAt   time.Time
}
