package http

import (
	"github.com/Royal17x/hireradar/internal/domain"
	logger "github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func (s *Server) handleRegistration(c *gin.Context) {

	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Request can't be read", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка прочтения запроса"})
		return
	}
	account, err := s.accountRepo.GetByEmail(c, req.Email)
	if err != nil {
		logger.Error("User not found", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}
	if account != nil {
		logger.Error("User already exists", "error", err)
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь уже существует"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Error hashing password", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хэширования пароля"})
		return
	}
	newAccount := domain.Account{
		Email:        req.Email,
		PasswordHash: string(hashed),
		CreatedAt:    time.Now(),
	}
	id, err := s.accountRepo.Save(c, newAccount)
	if err != nil {
		logger.Error("Error saving account", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения пользователя"})
		return
	}
	claims := &Claims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		logger.Error("Error signing token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подписи токена"})
		return
	}
	logger.Info("New account created", "email", req.Email)
	c.JSON(http.StatusCreated, gin.H{"token": tokenStr})
}

func (s *Server) handleLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Request can't be read", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка прочтения запроса"})
		return
	}
	if req.Email == "" || req.Password == "" {
		logger.Error("Incorrect credentials")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пустые поля почты или пароля"})
		return
	}
	account, err := s.accountRepo.GetByEmail(c, req.Email)
	if err != nil {
		logger.Error("Error getting account", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения аккаунта"})
		return
	}
	if account == nil {
		logger.Error("Account not found", "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Аккаунт не найден"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(req.Password))
	if err != nil {
		logger.Error("Incorrect password", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}
	claims := &Claims{
		UserID: account.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		logger.Error("Error signing token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подписи токена"})
		return
	}
	logger.Info("User successfully sign in", "email", req.Email)
	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

func (s *Server) handleVacancies(c *gin.Context) {
	userID := c.GetInt("user_id")
	vacancies, err := s.vacancyUcase.GetFiltered(c, userID)
	if err != nil {
		logger.Error("Failed to get vacancies", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить вакансии"})
		return
	}
	logger.Info("Vacancies are successfully shown")
	c.JSON(http.StatusOK, gin.H{"vacancies": vacancies})
}

func (s *Server) handleSetFilter(c *gin.Context) {
	userID := c.GetInt("user_id")
	var req struct {
		Keywords string `json:"keywords" binding:"required"`
		City     string `json:"city"`
		Grade    string `json:"grade"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Wrong request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong request"})
		return
	}
	filter := domain.Filter{
		UserID:    userID,
		Keywords:  req.Keywords,
		City:      req.City,
		Grade:     req.Grade,
		CreatedAt: time.Now(),
	}
	if err := s.filterRepo.Save(c, filter); err != nil {
		logger.Error("Error saving filter", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения фильтра"})
		return
	}
	logger.Info("Filter successfully saved")
	c.JSON(http.StatusCreated, gin.H{"filter": filter})
}

func (s *Server) handleGetFilters(c *gin.Context) {
	type filterResponse struct {
		ID       int    `json:"id"`
		Keywords string `json:"keywords"`
		City     string `json:"city"`
		Grade    string `json:"grade"`
		Active   bool   `json:"active"`
	}
	var response []filterResponse
	userID := c.GetInt("user_id")
	filters, err := s.filterRepo.GetByUserID(c, userID)
	if err != nil {
		logger.Error("Error getting filters", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting filters"})
		return
	}
	for i, f := range filters {
		response = append(response, filterResponse{
			ID:       f.ID,
			Keywords: f.Keywords,
			City:     f.City,
			Grade:    f.Grade,
			Active:   i == 0,
		})
	}
	logger.Info("Filters are successfully retrieved")
	c.JSON(http.StatusOK, gin.H{"filters": response})
}

func (s *Server) handleDeleteFilter(c *gin.Context) {
	var req struct {
		FilterID int `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Wrong request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong request"})
		return
	}
	if err := s.filterRepo.Delete(c, req.FilterID); err != nil {
		logger.Error("Error deleting filter", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting filter"})
		return
	}
	logger.Info("Filter successfully deleted")
	c.JSON(http.StatusOK, gin.H{"filter": nil})
}
