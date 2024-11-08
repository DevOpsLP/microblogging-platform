package persistence

import (
	"errors"
	"fmt"
	"time"

	"github.com/DevOpslp/microblogging-platform/user-service/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Método para encontrar un usuario dado un Username
func (repo *UserRepository) FindUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := repo.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("usuario no encontrado: %w", err)
		}
		return nil, err
	}
	return &user, nil
}

// Método para encontrar un usuario dado un ID
func (repo *UserRepository) FindUserByID(userID uint) (*domain.User, error) {
	var user domain.User
	if err := repo.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("usuario no encontrado: %w", err)
		}
		return nil, err
	}
	return &user, nil
}

// Método para obtener todos los usuarios con solo ID y Username
func (repo *UserRepository) GetAllUsers() ([]*domain.User, error) {
	var users []*domain.User
	if err := repo.db.Select("id", "username").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Métodos que usan user IDs
func (repo *UserRepository) FollowUser(userID, followID uint) error {
	// Buscar el usuario y el usuario a seguir usando los IDs
	var user, followUser domain.User
	if err := repo.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("no se encontró el usuario con ID %d: %w", userID, err)
	}
	if err := repo.db.First(&followUser, followID).Error; err != nil {
		return fmt.Errorf("no se encontró el usuario a seguir con ID %d: %w", followID, err)
	}

	// Intentar agregar la relación de seguimiento
	if err := repo.db.Model(&user).Association("Following").Append(&followUser); err != nil {
		return fmt.Errorf("no se pudo seguir al usuario con ID %d: %w", followID, err)
	}

	return nil
}

func (repo *UserRepository) UnfollowUser(userID, unfollowID uint) error {
	// Buscar el usuario y el usuario a dejar de seguir usando los IDs
	var user, unfollowUser domain.User
	if err := repo.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("no se encontró el usuario con ID %d: %w", userID, err)
	}
	if err := repo.db.First(&unfollowUser, unfollowID).Error; err != nil {
		return fmt.Errorf("no se encontró el usuario a dejar de seguir con ID %d: %w", unfollowID, err)
	}

	// Intentar eliminar la relación de seguimiento
	if err := repo.db.Model(&user).Association("Following").Delete(&unfollowUser); err != nil {
		return fmt.Errorf("no se pudo dejar de seguir al usuario con ID %d: %w", unfollowID, err)
	}

	return nil
}

func (repo *UserRepository) GetFollowers(userID uint) ([]*domain.User, error) {
	var user domain.User
	var followers []*domain.User

	// Buscar el usuario en la base de datos
	if err := repo.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Cargar la asociación de seguidores y asegurarse de que followers no sea nil
	if err := repo.db.Model(&user).Association("Followers").Find(&followers); err != nil {
		return nil, err
	}

	// Asegurar que followers esté inicializado, incluso si está vacío
	if followers == nil {
		followers = []*domain.User{}
	}

	return followers, nil
}

func (repo *UserRepository) GetFollowing(userID uint) ([]*domain.User, error) {
	var user domain.User
	var following []*domain.User

	// Buscar el usuario en la base de datos
	if err := repo.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Cargar la asociación de seguidos y asegurarse de que following no sea nil
	if err := repo.db.Model(&user).Association("Following").Find(&following); err != nil {
		return nil, err
	}

	// Asegurar que following esté inicializado, incluso si está vacío
	if following == nil {
		following = []*domain.User{}
	}

	return following, nil
}

var ErrUserAlreadyExists = errors.New("usuario ya registrado")

func (repo *UserRepository) RegisterUser(username, email string) (*domain.User, error) {
	// Verificar si el usuario ya existe
	var existingUser domain.User
	if err := repo.db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		return nil, ErrUserAlreadyExists
	}

	// Crear usuario con listas Following y Followers vacías
	user := domain.User{
		Username:  username,
		Email:     email,
		Following: []*domain.User{},
		Followers: []*domain.User{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := repo.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
