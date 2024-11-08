package persistence

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DevOpslp/microblogging-platform/user-service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() *gorm.DB {
	// Cargar las configuraciones desde las variables de entorno
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Construir el DSN (Data Source Name) usando las variables de entorno
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	// Migración automática del modelo User
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("Error al migrar el modelo User: %v", err)
	}

	// Verificar y crear datos iniciales
	seedUsers(db)

	return db
}

func seedUsers(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Count(&count)
	log.Printf("Número de usuarios en la base de datos: %d", count)

	if count == 0 {
		log.Println("Creando usuarios de ejemplo...")

		users := []domain.User{
			{Username: "user1", Email: "user1@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now(), Following: []*domain.User{}, Followers: []*domain.User{}},
			{Username: "user2", Email: "user2@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now(), Following: []*domain.User{}, Followers: []*domain.User{}},
			{Username: "user3", Email: "user3@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now(), Following: []*domain.User{}, Followers: []*domain.User{}},
			{Username: "user4", Email: "user4@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now(), Following: []*domain.User{}, Followers: []*domain.User{}},
		}

		// Insertar los usuarios en la base de datos
		if err := db.Create(&users).Error; err != nil {
			log.Fatalf("Error al crear usuarios de ejemplo: %v", err)
		}

		log.Println("Usuarios de ejemplo creados exitosamente.")
	} else {
		log.Println("Los usuarios de ejemplo ya existen en la base de datos.")
	}
}
