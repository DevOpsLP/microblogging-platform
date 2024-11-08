// internal/infrastructure/api/routes.go
package api

import (
	"github.com/DevOpslp/microblogging-platform/tweet-service/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, tweetRepo *persistence.TweetRepository) {
	handler := NewTweetHandler(tweetRepo)

	router.POST("/tweets", handler.CreateTweet)
	router.GET("/tweets", handler.GetAllTweets) // Nueva ruta para obtener todos los tweets
	router.GET("/tweets/:id", handler.GetTweet)
	router.GET("/tweets/user/:username", handler.GetTweetsByUser)
	router.DELETE("/tweets/:id", handler.DeleteTweet)

}
