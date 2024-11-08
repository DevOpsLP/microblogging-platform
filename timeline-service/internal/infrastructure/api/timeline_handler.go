package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/DevOpslp/microblogging-platform/timeline-service/internal/domain"
	"github.com/gin-gonic/gin"
)

type TimelineHandler struct {
	tweetServiceURL string
}

func NewTimelineHandler() *TimelineHandler {
	// Obtiene la URL de tweet-service desde las variables de entorno
	tweetServiceURL := os.Getenv("TWEET_SERVICE_URL")
	if tweetServiceURL == "" {
		tweetServiceURL = "http://localhost:8081/tweets" // Valor predeterminado para desarrollo local
	}
	return &TimelineHandler{tweetServiceURL: tweetServiceURL}
}

func (h *TimelineHandler) GetTimeline(c *gin.Context) {
	resp, err := http.Get(h.tweetServiceURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el timeline"})
		return
	}
	defer resp.Body.Close()

	var tweets []domain.Tweet
	if err := json.NewDecoder(resp.Body).Decode(&tweets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al decodificar tweets"})
		return
	}

	// Usamos MarshalIndent para embellecer el JSON
	prettyJSON, err := json.MarshalIndent(gin.H{"timeline": tweets}, "", "    ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al formatear el JSON"})
		return
	}

	// Establecemos el tipo de contenido como JSON y escribimos la respuesta
	c.Header("Content-Type", "application/json")
	c.Writer.Write(prettyJSON)
}
