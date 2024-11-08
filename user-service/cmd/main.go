package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/DevOpslp/microblogging-platform/user-service/internal/infrastructure/api"
	"github.com/DevOpslp/microblogging-platform/user-service/internal/infrastructure/persistence"
)

func main() {
	// Configuración de la base de datos
	db := persistence.NewDB()

	// Configuración del repositorio de usuarios
	userRepository := persistence.NewUserRepository(db)

	// Inicia el enrutador de Gin
	router := gin.Default()

	// Pasar userRepository a SetupRoutes
	api.SetupRoutes(router, *userRepository)

	// Obtiene el puerto desde las variables de entorno o usa 8080 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Inicia el servidor en el puerto especificado
	log.Printf("Iniciando el servidor en el puerto %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %v", err)
	}
}
