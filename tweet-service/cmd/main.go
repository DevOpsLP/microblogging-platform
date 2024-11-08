package main

import (
	"log"
	"os"

	"github.com/DevOpslp/microblogging-platform/tweet-service/internal/domain"
	"github.com/DevOpslp/microblogging-platform/tweet-service/internal/infrastructure/api"
	"github.com/DevOpslp/microblogging-platform/tweet-service/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Obtener las variables de entorno para la conexión a la base de datos
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	userServiceURL := os.Getenv("USER_SERVICE_URL")

	// Construir el DSN (Data Source Name) usando las variables de entorno
	tweetDBDSN := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable"

	tweetDB, err := gorm.Open(postgres.Open(tweetDBDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a tweetdb: %v", err)
	}

	if err := tweetDB.AutoMigrate(&domain.Tweet{}); err != nil {
		log.Fatalf("Error al migrar el modelo Tweet: %v", err)
	}

	// Crear los repositorios de usuario y tweet
	userRepo := persistence.NewHTTPUserRepository(userServiceURL) // URL de `user-service`
	tweetRepo := persistence.NewTweetRepository(tweetDB, userRepo)

	// Iniciar el servidor HTTP
	router := gin.Default()
	api.SetupRoutes(router, tweetRepo)

	// Escuchar en el puerto 8081
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Error al iniciar el servidor de Tweet-Service: %v", err)
	}
}
