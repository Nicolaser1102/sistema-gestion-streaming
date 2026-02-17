package handlers

import (
	"net/http"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"

	"streaming-system/internal/models"
	"streaming-system/internal/repositories"
)

type ContentHandler struct {
	contentRepo *repositories.ContentRepoJSON
}

func NewContentHandler(repo *repositories.ContentRepoJSON) *ContentHandler {
	return &ContentHandler{contentRepo: repo}
}

// GET /api/contents
func (h *ContentHandler) ListContents(c *gin.Context) {

	query := c.Query("q")
	genre := c.Query("genre")
	ctype := c.Query("type")
	yearStr := c.Query("year")

	var year int
	if yearStr != "" {
		year, _ = strconv.Atoi(yearStr)
	}

	contents, err := h.contentRepo.SearchAndFilter(query, genre, ctype, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}

	c.JSON(http.StatusOK, contents)
}

// GET /api/contents/:id
func (h *ContentHandler) GetContentByID(c *gin.Context) {

	id := c.Param("id")

	content, err := h.contentRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}

	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contenido no encontrado"})
		return
	}

	c.JSON(http.StatusOK, content)

}

// GET /api/admin/contents
func (h *ContentHandler) AdminListContents(c *gin.Context) {
	contents, err := h.contentRepo.GetAllIncludingInactive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}
	c.JSON(http.StatusOK, contents)
}

// POST /api/admin/contents
func (h *ContentHandler) CreateContent(c *gin.Context) {

	var content models.Content

	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	content.ID = generateID()
	content.Active = true
	content.CreatedAt = time.Now()
	content.UpdatedAt = time.Now()

	if err := h.contentRepo.Create(content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear"})
		return
	}

	c.JSON(http.StatusCreated, content)
}

// PUT /api/admin/contents/:id
func (h *ContentHandler) UpdateContent(c *gin.Context) {

	id := c.Param("id")

	var updated models.Content

	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if err := h.contentRepo.Update(id, updated); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contenido no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contenido actualizado"})
}

// DELETE /api/admin/contents/:id
func (h *ContentHandler) DeleteContent(c *gin.Context) {

	id := c.Param("id")

	if err := h.contentRepo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contenido no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contenido desactivado"})
}
