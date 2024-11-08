package persistence

import (
	"fmt"
	"log"
	"testing"

	"github.com/DevOpslp/microblogging-platform/tweet-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

const (
	userDBDSN  = "host=localhost user=devuser password=devpassword dbname=userdb port=5432 sslmode=disable"
	tweetDBDSN = "host=localhost user=devuser password=devpassword dbname=tweetdb port=5432 sslmode=disable"
)

func TestDBConnections(t *testing.T) {
	// Conectar a userdb
	userDB := NewDB(userDBDSN)
	defer func() {
		sqlDB, _ := userDB.DB()
		sqlDB.Close()
	}()

	assert.NotNil(t, userDB, "La conexión a la base de datos userdb debería haberse establecido")
	fmt.Println("Conexión exitosa a userdb")

	// Conectar a tweetdb
	tweetDB := NewDB(tweetDBDSN)
	defer func() {
		sqlDB, _ := tweetDB.DB()
		sqlDB.Close()
	}()

	assert.NotNil(t, tweetDB, "La conexión a la base de datos tweetdb debería haberse establecido")
	fmt.Println("Conexión exitosa a tweetdb")
}

func TestAutoMigration(t *testing.T) {
	// Conectar a tweetDB y realizar migración
	tweetDB := NewDB(tweetDBDSN)
	defer func() {
		sqlDB, _ := tweetDB.DB()
		sqlDB.Close()
	}()

	err := tweetDB.AutoMigrate(&domain.Tweet{})

	if err := tweetDB.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("Error al migrar el modelo User: %v", err)
	}

	if err := tweetDB.AutoMigrate(&domain.Tweet{}); err != nil {
		log.Fatalf("Error al migrar el modelo Tweet: %v", err)
	}

	assert.NoError(t, err, "La migración automática de las tablas debería completarse sin errores")
}
