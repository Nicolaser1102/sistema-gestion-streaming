package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"streaming-system/internal/repositories"
)

type MyListHandler struct {
	myListRepo  *repositories.MyListRepoJSON
	contentRepo *repositories.ContentRepoJSON
}

func NewMyListHandler(myListRepo *repositories.MyListRepoJSON, contentRepo *repositories.ContentRepoJSON) *MyListHandler {
	return &MyListHandler{
		myListRepo:  myListRepo,
		contentRepo: contentRepo,
	}
}

// GET /api/my-list
func (h *MyListHandler) GetMyList(c *gin.Context) {
	userIDVal, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
		return
	}
	userID, _ := userIDVal.(string)

	items, err := h.myListRepo.GetByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}

	// devolver contenidos (más útil para frontend)
	var contents []any
	for _, it := range items {
		content, err := h.contentRepo.FindByID(it.ContentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
			return
		}
		if content != nil {
			contents = append(contents, content)
		}
	}

	c.JSON(http.StatusOK, contents)
}

// POST /api/my-list/:contentId
func (h *MyListHandler) AddToMyList(c *gin.Context) {
	userIDVal, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
		return
	}
	userID, _ := userIDVal.(string)

	contentID := c.Param("contentId")

	// validar que exista el contenido
	content, err := h.contentRepo.FindByID(contentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
		return
	}
	if content == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contenido no encontrado"})
		return
	}

	if err := h.myListRepo.Add(userID, contentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo agregar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agregado a Mi lista"})
}

// DELETE /api/my-list/:contentId
func (h *MyListHandler) RemoveFromMyList(c *gin.Context) {
	userIDVal, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autenticado"})
		return
	}
	userID, _ := userIDVal.(string)

	contentID := c.Param("contentId")

	if err := h.myListRepo.Remove(userID, contentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Eliminado de Mi lista"})
}
