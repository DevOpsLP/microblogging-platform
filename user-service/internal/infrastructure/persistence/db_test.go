package persistence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	db := NewDB()

	// Verificar que la conexión a la base de datos no sea nil
	assert.NotNil(t, db, "La conexión a la base de datos debería haber sido establecida")

	// Verificar que se puede hacer un ping a la base de datos para confirmar la conexión
	sqlDB, err := db.DB()
	assert.NoError(t, err, "Error al obtener el objeto DB de GORM")

	err = sqlDB.Ping()
	assert.NoError(t, err, "La conexión a la base de datos no fue exitosa")
}
