package persistence

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DevOpslp/microblogging-platform/tweet-service/internal/domain"
)

type HTTPUserRepository struct {
	baseURL string
}

type UserRepository interface {
	FindUserByUsername(username string) (*domain.User, error)
	FindUserByID(userID uint) (*domain.User, error) // Añadir este método
}

func NewHTTPUserRepository(baseURL string) *HTTPUserRepository {
	return &HTTPUserRepository{baseURL: baseURL}
}

func (repo *HTTPUserRepository) FindUserByUsername(username string) (*domain.User, error) {
	url := fmt.Sprintf("%s/user/%s", repo.baseURL, username) // Asegúrate de usar /user/:username
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: user-service devolvió estado %d", resp.StatusCode)
	}

	var result struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &domain.User{ID: result.UserID, Username: result.Username}, nil
}

func (repo *HTTPUserRepository) FindUserByID(userID uint) (*domain.User, error) {
	url := fmt.Sprintf("%s/user-by-id/%d", repo.baseURL, userID) // Endpoint nuevo
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: user-service devolvió estado %d", resp.StatusCode)
	}

	var result struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &domain.User{ID: result.UserID, Username: result.Username}, nil
}
