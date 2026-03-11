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
	favoriteRepo domain.FavoriteRepository
	parserQuery  string
	jwtSecret    string
}

func NewServer(vacancyUcase *usecase.VacancyUsecase, userRepo domain.UserRepository, filterRepo domain.FilterRepository, accountRepo domain.AccountRepository, favoriteRepo domain.FavoriteRepository, parserQuery string, jwtSecret string) *Server {
	s := &Server{
		router:       gin.Default(),
		vacancyUcase: vacancyUcase,
		userRepo:     userRepo,
		filterRepo:   filterRepo,
		accountRepo:  accountRepo,
		favoriteRepo: favoriteRepo,
		parserQuery:  parserQuery,
		jwtSecret:    jwtSecret,
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	// auth
	api := s.router.Group("/api")
	api.POST("/register", s.handleRegistration)
	api.POST("/login", s.handleLogin)

	auth := api.Group("/")
	auth.Use(AuthMiddleware(s.jwtSecret))
	// vacancies
	auth.GET("/vacancies", s.handleGetVacancies)
	auth.GET("/vacancies/:hh_id", s.handleGetVacancy)
	auth.POST("/vacancies/refresh", s.handleRefreshVacancies)
	// filters
	auth.POST("/filters", s.handleSetFilter)
	auth.GET("/filter/:userID", s.handleGetFilters)
	auth.DELETE("/filters/:id", s.handleDeleteFilter)
	// favorites
	auth.POST("/favorites", s.handleAddFavorite)
	auth.GET("/favorites", s.handleGetFavorites)
	auth.DELETE("/favorites/:hh_id", s.handleDeleteFavorite)
	// features
	auth.GET("/profile", s.handleProfile)
	auth.GET("/stats", s.handleStats)
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
