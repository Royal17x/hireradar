package hh

import (
	"context"
	"encoding/json"
	"github.com/Royal17x/hireradar/internal/domain"
	"net/http"
	"time"
)

type Client struct {
	httpclient *http.Client
	baseURL    string
}

func New() *Client {
	return &Client{
		httpclient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.hh.ru",
	}
}

type hhResponse struct {
	Items []hhVacancy `json:"items"`
}

type hhVacancy struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Employer     hhEmployer `json:"employer"`
	Salary       *hhSalary  `json:"salary"`
	AlternateURL string     `json:"alternate_url"`
	PublishedAt  string     `json:"published_at"`
	CreatedAt    string     `json:"created_at"`
}

type hhEmployer struct {
	Name string `json:"name"`
}

type hhSalary struct {
	From int `json:"from"`
	To   int `json:"to"`
}

func (c *Client) FetchVacancies(ctx context.Context, query string) ([]domain.Vacancy, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/vacancies", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("text", query)
	q.Add("area", "1")
	q.Add("per_page", "100")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response hhResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	var vacancies []domain.Vacancy
	for _, vacancy := range response.Items {
		t1, err := time.Parse(time.RFC3339, vacancy.CreatedAt)
		if err != nil {
			return nil, err
		}
		t2, err := time.Parse(time.RFC3339, vacancy.PublishedAt)
		if err != nil {
			return nil, err
		}
		v := domain.Vacancy{
			HhID:        vacancy.ID,
			Title:       vacancy.Name,
			Company:     vacancy.Employer.Name,
			URL:         vacancy.AlternateURL,
			PublishedAt: t1,
			CreatedAt:   t2,
		}
		if vacancy.Salary != nil {
			from := vacancy.Salary.From
			to := vacancy.Salary.To
			v.SalaryFrom = &from
			v.SalaryTo = &to
		}
		vacancies = append(vacancies, v)
	}
	return vacancies, nil
}
