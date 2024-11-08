package api

import (
	"github.com/DevOpslp/microblogging-platform/user-service/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userRepo persistence.UserRepository) {
	handler := NewUserHandler(userRepo)

	router.POST("/follow", handler.FollowUser)
	router.POST("/unfollow", handler.UnfollowUser)
	router.POST("/register", handler.RegisterUser)
	router.GET("/followers", handler.GetFollowers)
	router.GET("/following", handler.GetFollowing)
	router.GET("/users", handler.GetAllUsers)
	router.GET("/user/:username", handler.GetUserByUsername)
	router.GET("/user-by-id/:id", handler.GetUserByID)
}
