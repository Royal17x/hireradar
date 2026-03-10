package http

import (
	"github.com/Royal17x/hireradar/internal/domain"
	"github.com/Royal17x/hireradar/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router       *gin.Engine
	vacancyUcase *usecase.VacancyUsecase
	userRepo     domain.UserRepository
	filterRepo   domain.FilterRepository
	accountRepo  domain.AccountRepository
	jwtSecret    string
}

func NewServer(vacancyUcase *usecase.VacancyUsecase, userRepo domain.UserRepository, filterRepo domain.FilterRepository, accountRepo domain.AccountRepository, jwtSecret string) *Server {
	s := &Server{
		router:       gin.Default(),
		vacancyUcase: vacancyUcase,
		userRepo:     userRepo,
		filterRepo:   filterRepo,
		accountRepo:  accountRepo,
		jwtSecret:    jwtSecret,
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	api := s.router.Group("/api")
	api.POST("/register", s.handleRegistration)
	api.POST("/login", s.handleLogin)

	auth := api.Group("/")
	auth.Use(AuthMiddleware(s.jwtSecret))
	auth.GET("/vacancies", s.handleVacancies)
	auth.POST("/filters", s.handleSetFilter)
	auth.GET("/filter/:userID", s.handleGetFilters)
	auth.DELETE("/filters/:id", s.handleDeleteFilter)

}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
