// internal/infrastructure/api/tweet_handler.go
package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/DevOpslp/microblogging-platform/tweet-service/internal/domain"
	"github.com/DevOpslp/microblogging-platform/tweet-service/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
)

type TweetHandler struct {
	repo *persistence.TweetRepository
}

type TweetResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func formatTweetResponse(tweet domain.TweetWithUser) TweetResponse {
	return TweetResponse{
		ID:        tweet.ID,
		Username:  tweet.Username,
		Content:   tweet.Content,
		CreatedAt: tweet.CreatedAt.Format(time.RFC3339),
		UpdatedAt: tweet.UpdatedAt.Format(time.RFC3339),
	}
}

func NewTweetHandler(repo *persistence.TweetRepository) *TweetHandler {
	return &TweetHandler{repo: repo}
}

func (h *TweetHandler) CreateTweet(c *gin.Context) {
	var body struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contenido del tweet inválido"})
		return
	}

	username := c.GetHeader("Username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username no proporcionado en el header"})
		return
	}

	tweet, err := h.repo.CreateTweet(username, body.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el tweet"})
		return
	}

	c.JSON(http.StatusCreated, tweet)
}

func (h *TweetHandler) GetTweet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tweet inválido"})
		return
	}

	tweet, err := h.repo.GetTweetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tweet no encontrado"})
		return
	}

	c.JSON(http.StatusOK, tweet)
}

func (h *TweetHandler) GetAllTweets(c *gin.Context) {
	tweets, err := h.repo.GetAllTweets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los tweets"})
		return
	}

	// Formatear cada tweet para la respuesta
	var response []TweetResponse
	for _, tweet := range tweets {
		response = append(response, formatTweetResponse(tweet))
	}

	c.JSON(http.StatusOK, response)
}

func (h *TweetHandler) GetTweetsByUser(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username inválido"})
		return
	}

	tweets, err := h.repo.GetTweetsByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los tweets"})
		return
	}

	c.JSON(http.StatusOK, tweets)
}

func (h *TweetHandler) DeleteTweet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tweet inválido"})
		return
	}

	if err := h.repo.DeleteTweetByID(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar el tweet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tweet eliminado exitosamente"})
}
