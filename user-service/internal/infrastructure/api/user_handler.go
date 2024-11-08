package api

import (
	"net/http"
	"strconv"

	"github.com/DevOpslp/microblogging-platform/user-service/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo persistence.UserRepository
}

func NewUserHandler(userRepo persistence.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var body struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de registro inválidos"})
		return
	}

	user, err := h.userRepo.RegisterUser(body.Username, body.Email)
	if err != nil {
		if err == persistence.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Usuario ya registrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el usuario"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Usuario registrado exitosamente",
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
		return
	}

	user, err := h.userRepo.FindUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": user.ID, "username": user.Username})
}

// Obtener User-ID a partir del username en el header
func (h *UserHandler) getUserIDFromUsernameHeader(c *gin.Context) (uint, error) {
	username := c.GetHeader("Username")
	if username == "" {
		return 0, http.ErrNoCookie
	}

	user, err := h.userRepo.FindUserByUsername(username)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := h.userRepo.FindUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": user.ID, "username": user.Username})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los usuarios"})
		return
	}

	// Responder solo con ID y Username
	usersResponse := make([]gin.H, 0)
	for _, user := range users {
		usersResponse = append(usersResponse, gin.H{
			"user_id":  user.ID,
			"username": user.Username,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": usersResponse})
}

func (h *UserHandler) FollowUser(c *gin.Context) {
	userID, err := h.getUserIDFromUsernameHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username no válido o no encontrado"})
		return
	}

	var body struct {
		FollowUsername string `json:"follow_username"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.FollowUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "FollowUsername no válido"})
		return
	}

	// Encontrar el usuario a seguir usando el username
	followUser, err := h.userRepo.FindUserByUsername(body.FollowUsername)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario a seguir no encontrado"})
		return
	}

	// Ahora pasamos los IDs
	if err := h.userRepo.FollowUser(userID, followUser.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo seguir al usuario"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuario seguido exitosamente"})
}

func (h *UserHandler) UnfollowUser(c *gin.Context) {
	userID, err := h.getUserIDFromUsernameHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username no válido o no encontrado"})
		return
	}

	var body struct {
		UnfollowUsername string `json:"unfollow_username"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.UnfollowUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UnfollowUsername no válido"})
		return
	}

	unfollowUser, err := h.userRepo.FindUserByUsername(body.UnfollowUsername)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario a dejar de seguir no encontrado"})
		return
	}

	if err := h.userRepo.UnfollowUser(userID, unfollowUser.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo dejar de seguir al usuario"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuario dejado de seguir exitosamente"})
}

func (h *UserHandler) GetFollowers(c *gin.Context) {
	userID, err := h.getUserIDFromUsernameHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username no válido o no encontrado"})
		return
	}

	followers, err := h.userRepo.GetFollowers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener la lista de seguidores"})
		return
	}

	followersResponse := make([]gin.H, 0)
	for _, follower := range followers {
		followersResponse = append(followersResponse, gin.H{"username": follower.Username})
	}

	c.JSON(http.StatusOK, gin.H{"followers": followersResponse})
}

func (h *UserHandler) GetFollowing(c *gin.Context) {
	userID, err := h.getUserIDFromUsernameHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username no válido o no encontrado"})
		return
	}

	following, err := h.userRepo.GetFollowing(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener la lista de usuarios seguidos"})
		return
	}

	followingResponse := make([]gin.H, 0)
	for _, user := range following {
		followingResponse = append(followingResponse, gin.H{"username": user.Username})
	}

	c.JSON(http.StatusOK, gin.H{"following": followingResponse})
}
