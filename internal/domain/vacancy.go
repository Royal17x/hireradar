package domain

import "time"

type Vacancy struct {
	VacancyID   int       `json:"-"`
	HhID        string    `json:"hh_id"`
	Title       string    `json:"title"`
	City        string    `json:"city"`
	Company     string    `json:"company"`
	URL         string    `json:"url"`
	SalaryFrom  *int      `json:"salary_from"`
	SalaryTo    *int      `json:"salary_to"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
}
