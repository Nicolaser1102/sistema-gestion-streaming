package handlers

import (
	"net/http"

	"streaming-system/internal/models"
	"streaming-system/internal/repositories"

	"github.com/gin-gonic/gin"
)

type PlaybackHandler struct {
	progressRepo *repositories.ProgressRepoJSON
	contentRepo  *repositories.ContentRepoJSON
}

func NewPlaybackHandler(progressRepo *repositories.ProgressRepoJSON, contentRepo *repositories.ContentRepoJSON) *PlaybackHandler {
	return &PlaybackHandler{progressRepo: progressRepo, contentRepo: contentRepo}
}

// GET /api/playback/:contentId  -> devuelve progreso actual
func (h *PlaybackHandler) GetProgress(c *gin.Context) {
	userIDVal, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
		return
	}
	userID, _ := userIDVal.(string)

	contentID := c.Param("contentId")

	// validar contenido existe
	content, err := h.contentRepo.FindByID(contentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}
	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contenido no encontrado"})
		return
	}

	p, err := h.progressRepo.GetByUserAndContent(userID, contentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}

	if p == nil {
		c.JSON(http.StatusOK, gin.H{"seconds": 0, "percent": 0, "completed": false})
		return
	}

	c.JSON(http.StatusOK, p)
}

type UpdateProgressRequest struct {
	Seconds int     `json:"seconds"`
	Percent float64 `json:"percent"`
}

// PUT /api/playback/:contentId/progress
func (h *PlaybackHandler) UpdateProgress(c *gin.Context) {
	userIDVal, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
		return
	}
	userID, _ := userIDVal.(string)

	contentID := c.Param("contentId")

	// validar contenido
	content, err := h.contentRepo.FindByID(contentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}
	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contenido no encontrado"})
		return
	}

	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if req.Seconds < 0 {
		req.Seconds = 0
	}
	if req.Percent < 0 {
		req.Percent = 0
	}
	if req.Percent > 100 {
		req.Percent = 100
	}

	completed := req.Percent >= 90 // regla académica: >=90% se marca completado

	p := models.PlaybackProgress{
		UserID:    userID,
		ContentID: contentID,
		Seconds:   req.Seconds,
		Percent:   req.Percent,
		Completed: completed,
	}

	if err := h.progressRepo.Upsert(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar progreso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Progreso guardado", "completed": completed})
}
