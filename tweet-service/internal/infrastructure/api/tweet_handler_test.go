// internal/infrastructure/api/tweet_handler_test.go
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DevOpslp/microblogging-platform/tweet-service/internal/domain"
	"github.com/DevOpslp/microblogging-platform/tweet-service/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func setupTestDB() {
	dsn := "host=localhost user=devuser password=devpassword dbname=tweetdb port=5432 sslmode=disable"
	var err error
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos de pruebas: %v", err)
	}

	// Migración automática de la base de datos para el modelo Tweet
	testDB.AutoMigrate(&domain.Tweet{})
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	userRepo := persistence.NewHTTPUserRepository("http://localhost:8080")
	tweetRepo := persistence.NewTweetRepository(testDB, userRepo)
	router := gin.Default()
	SetupRoutes(router, tweetRepo)
	return router
}

// obtener un usuario aleatorio de user-service
func getRandomUser() (*domain.User, error) {
	resp, err := http.Get("http://localhost:8080/users")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: user-service devolvió estado %d", resp.StatusCode)
	}

	// Decodificar la respuesta JSON en una estructura que contiene una lista de usuarios
	var response struct {
		Users []domain.User `json:"users"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Verificar que haya al menos un usuario en la lista
	if len(response.Users) == 0 {
		return nil, fmt.Errorf("no se encontraron usuarios en user-service")
	}

	// Seleccionar un usuario aleatorio de la lista
	return &response.Users[0], nil // Puedes ajustar aquí si quieres un usuario aleatorio diferente
}

func TestTweetFlowWithHTTPUserRepo(t *testing.T) {
	setupTestDB()

	// Obtener dos usuarios aleatorios desde user-service
	user1, err := getRandomUser()
	assert.NoError(t, err, "Debe haber al menos un usuario en user-service para realizar la prueba")

	user2, err := getRandomUser()
	assert.NoError(t, err, "Debe haber al menos un segundo usuario en user-service para realizar la prueba")

	router := setupTestRouter()

	// Crear un tweet para user1
	t.Run("Crear Tweet para User1", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"content":"Tweet de prueba para user1"}`
		req, _ := http.NewRequest("POST", "/tweets", bytes.NewBufferString(reqBody))
		req.Header.Set("Username", user1.Username)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"Tweet de prueba para user1"`)
	})

	// Crear un tweet para user2
	t.Run("Crear Tweet para User2", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"content":"Tweet de prueba para user2"}`
		req, _ := http.NewRequest("POST", "/tweets", bytes.NewBufferString(reqBody))
		req.Header.Set("Username", user2.Username)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"Tweet de prueba para user2"`)
	})

	// Obtener tweets de user1
	t.Run("Obtener Tweets de User1", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tweets/user/"+user1.Username, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), user1.Username)
	})

	// Obtener todos los tweets
	t.Run("Obtener Todos los Tweets", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tweets", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Tweet de prueba para user1")
		assert.Contains(t, w.Body.String(), "Tweet de prueba para user2")
	})

	// Eliminar tweet de user1
	t.Run("Eliminar Tweet de User1", func(t *testing.T) {
		// Obtener ID de un tweet de user1
		tweetID := uint(1) // Reemplaza con un ID válido o realiza un GET para obtener un tweet específico.
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/tweets/%d", tweetID), nil)
		req.Header.Set("Username", user1.Username)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Tweet eliminado exitosamente")
	})
}
