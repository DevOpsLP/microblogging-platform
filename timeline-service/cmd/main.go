package main

import (
	"log"

	"github.com/DevOpslp/microblogging-platform/timeline-service/internal/infrastructure/api"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Crear instancia de TimelineHandler
	timelineHandler := api.NewTimelineHandler()

	// Configurar rutas con la instancia de handler
	api.SetupRoutes(router, timelineHandler)

	// Iniciar el servidor en el puerto 8082
	log.Fatal(router.Run(":8082"))
}
