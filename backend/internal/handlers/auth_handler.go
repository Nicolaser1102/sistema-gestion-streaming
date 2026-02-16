package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"streaming-system/internal/models"
	"streaming-system/internal/repositories"
	"streaming-system/internal/utils"
)

type AuthHandler struct {
	userRepo   *repositories.UserRepoJSON
	jwtSecret  string
	jwtExpires time.Duration
}

func NewAuthHandler(userRepo *repositories.UserRepoJSON, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		jwtExpires: 24 * time.Hour, // token válido 24h (puedes cambiar)
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name es requerido"})
		return
	}
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email inválido"})
		return
	}
	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password mínimo 6 caracteres"})
		return
	}

	existingUser, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "El email ya está registrado"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encriptando password"})
		return
	}

	user := models.User{
		ID:           generateID(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleUser,
		CreatedAt:    time.Now(),
	}

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar usuario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || !strings.Contains(req.Email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email inválido"})
		return
	}
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password es requerido"})
		return
	}

	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	// comparar hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
		return
	}

	// generar token JWT
	token, err := utils.GenerateToken(h.jwtSecret, user.ID, user.Role, h.jwtExpires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func generateID() string {
	// Para la clase sirve, luego podemos mejorarlo a UUID.
	return time.Now().Format("20060102150405")
}

// GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userIDVal, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
		return
	}
	userID, _ := userIDVal.(string)

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no existe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"name":      user.Name,
		"email":     user.Email,
		"role":      user.Role,
		"createdAt": user.CreatedAt,
	})
}
