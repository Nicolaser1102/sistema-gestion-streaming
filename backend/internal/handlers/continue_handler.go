package handlers

import (
	"net/http"

	"streaming-system/internal/repositories"

	"github.com/gin-gonic/gin"
)

type ContinueHandler struct {
	progressRepo *repositories.ProgressRepoJSON
	contentRepo  *repositories.ContentRepoJSON
}

func NewContinueHandler(p *repositories.ProgressRepoJSON, c *repositories.ContentRepoJSON) *ContinueHandler {
	return &ContinueHandler{progressRepo: p, contentRepo: c}
}

func (h *ContinueHandler) GetContinueWatching(c *gin.Context) {
	userIDVal, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
		return
	}
	userID := userIDVal.(string)

	progressList, err := h.progressRepo.GetByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}

	// devolver contenido + progreso
	type Item struct {
		Content interface{} `json:"content"`
		Seconds int         `json:"seconds"`
		Percent float64     `json:"percent"`
	}

	out := []Item{}
	for _, p := range progressList {
		content, err := h.contentRepo.FindByID(p.ContentID)
		if err != nil {
			continue
		}
		if content == nil {
			continue
		}
		out = append(out, Item{
			Content: content,
			Seconds: p.Seconds,
			Percent: p.Percent,
		})
	}

	c.JSON(http.StatusOK, out)
}
