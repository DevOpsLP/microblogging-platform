package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DevOpslp/microblogging-platform/user-service/internal/domain"
	"github.com/DevOpslp/microblogging-platform/user-service/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var userRepo *persistence.UserRepository

// setupTestDB configura la base de datos de prueba y migra el esquema necesario
func setupTestDB() {
	dsn := "host=localhost user=devuser password=devpassword dbname=userdb port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("No se pudo conectar a la base de datos de pruebas")
	}

	// Migrar el esquema y crear el repositorio
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		panic("No se pudo migrar el esquema de User")
	}
	userRepo = persistence.NewUserRepository(db)
}

// cleanDatabase elimina usuarios específicos y sus relaciones para mantener la base de datos limpia después de las pruebas
func cleanDatabase(userIDs ...uint) {
	// Eliminar relaciones de seguidores/seguidores
	for _, id := range userIDs {
		db.Exec("DELETE FROM user_followers WHERE follower_id = ? OR user_id = ?", id, id)
	}
	// Eliminar usuarios
	for _, id := range userIDs {
		db.Exec("DELETE FROM users WHERE id = ?", id)
	}
}

// setupTestRouter configura el router de Gin para las pruebas
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	SetupRoutes(router, *userRepo)
	return router
}

// generateRandomUser crea un usuario con un nombre de usuario y email aleatorios para las pruebas
func generateRandomUser() (*domain.User, error) {
	username := fmt.Sprintf("test_user%d", rand.Intn(1000000))
	email := fmt.Sprintf("%s@example.com", username)
	user := domain.User{
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// TestUserFollowUnfollowFlow prueba el flujo completo de seguir y dejar de seguir usuarios
func TestUserFollowUnfollowFlow(t *testing.T) {
	// Inicializar la semilla aleatoria
	rand.Seed(time.Now().UnixNano())

	setupTestDB()

	// Crear usuarios de prueba
	user1, err := generateRandomUser()
	assert.NoError(t, err, "No se pudo crear el usuario1 en userDB")
	user2, err := generateRandomUser()
	assert.NoError(t, err, "No se pudo crear el usuario2 en userDB")

	// Limpiar la base de datos al finalizar la prueba
	defer cleanDatabase(user1.ID, user2.ID)

	router := setupTestRouter()

	// Paso 1: Usuario 1 sigue a Usuario 2
	t.Run("Usuario1 sigue a Usuario2", func(t *testing.T) {
		followRequest := fmt.Sprintf(`{"follow_username": "%s"}`, user2.Username)
		req, _ := http.NewRequest("POST", "/follow", strings.NewReader(followRequest))
		req.Header.Set("Username", user1.Username) // Usar el Username en el header
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "Usuario seguido exitosamente"}`, w.Body.String())
	})

	// Paso 2: Verificar que Usuario 1 está siguiendo a Usuario 2
	t.Run("Verificar que Usuario1 está siguiendo a Usuario2", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/following", nil)
		req.Header.Set("Username", user1.Username)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expectedFollowingResponse := fmt.Sprintf(`{"following": [{"username": "%s"}]}`, user2.Username)
		assert.JSONEq(t, expectedFollowingResponse, w.Body.String())
	})

	// Paso 3: Verificar que Usuario 2 tiene a Usuario 1 como seguidor
	t.Run("Verificar que Usuario2 tiene a Usuario1 como seguidor", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/followers", nil)
		req.Header.Set("Username", user2.Username)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expectedFollowersResponse := fmt.Sprintf(`{"followers": [{"username": "%s"}]}`, user1.Username)
		assert.JSONEq(t, expectedFollowersResponse, w.Body.String())
	})

	// Paso 4: Usuario 1 deja de seguir a Usuario 2
	t.Run("Usuario1 deja de seguir a Usuario2", func(t *testing.T) {
		unfollowRequest := fmt.Sprintf(`{"unfollow_username": "%s"}`, user2.Username)
		req, _ := http.NewRequest("POST", "/unfollow", strings.NewReader(unfollowRequest))
		req.Header.Set("Username", user1.Username) // Usar el Username en el header
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "Usuario dejado de seguir exitosamente"}`, w.Body.String())
	})

	// Paso 5: Confirmar que la lista de "following" de Usuario 1 está vacía
	t.Run("Verificar que la lista de 'following' de Usuario1 está vacía", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/following", nil)
		req.Header.Set("Username", user1.Username)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expectedEmptyFollowingResponse := `{"following": []}`
		assert.JSONEq(t, expectedEmptyFollowingResponse, w.Body.String())
	})

	// Paso 6: Confirmar que la lista de "followers" de Usuario 2 está vacía
	t.Run("Verificar que la lista de 'followers' de Usuario2 está vacía", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/followers", nil)
		req.Header.Set("Username", user2.Username)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expectedEmptyFollowersResponse := `{"followers": []}`
		assert.JSONEq(t, expectedEmptyFollowersResponse, w.Body.String())
	})
}
