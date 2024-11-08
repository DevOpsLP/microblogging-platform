// internal/infrastructure/ersistence/tweet_repository.go
package persistence

import (
	"errors"
	"fmt"
	"time"

	"github.com/DevOpslp/microblogging-platform/tweet-service/internal/domain"
	"gorm.io/gorm"
)

type TweetRepository struct {
	tweetDB  *gorm.DB
	userRepo UserRepository
}

func NewTweetRepository(tweetDB *gorm.DB, userRepo UserRepository) *TweetRepository {
	return &TweetRepository{tweetDB: tweetDB, userRepo: userRepo}
}

// Crear un tweet usando el username del `user-service`
func (repo *TweetRepository) CreateTweet(username, content string) (*domain.Tweet, error) {
	// Obtener `UserID` desde `user-service`
	user, err := repo.userRepo.FindUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado en user-service: %w", err)
	}

	// Crear el tweet asociado a `UserID`
	tweet := &domain.Tweet{
		UserID:    user.ID,
		Content:   content,
		CreatedAt: time.Now(), // Asignar explícitamente la fecha de creación

	}
	if err := repo.tweetDB.Create(tweet).Error; err != nil {
		return nil, err
	}

	return tweet, nil
}

// Obtener tweets por username
func (repo *TweetRepository) GetTweetsByUsername(username string) ([]domain.Tweet, error) {
	user, err := repo.userRepo.FindUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	var tweets []domain.Tweet
	if err := repo.tweetDB.Where("user_id = ?", user.ID).Find(&tweets).Error; err != nil {
		return nil, err
	}
	return tweets, nil
}

// Obtener todos los tweets
// Obtener todos los tweets con información del usuario
func (repo *TweetRepository) GetAllTweets() ([]domain.TweetWithUser, error) {
	var tweets []domain.Tweet
	if err := repo.tweetDB.Find(&tweets).Error; err != nil {
		fmt.Printf("Error al obtener tweets: %v\n", err) // Agrega el registro
		return nil, fmt.Errorf("error al obtener tweets: %w", err)
	}

	var tweetsWithUser []domain.TweetWithUser
	for _, tweet := range tweets {
		user, err := repo.userRepo.FindUserByID(tweet.UserID)
		if err != nil {
			fmt.Printf("Error al obtener usuario para tweet %d: %v\n", tweet.ID, err) // Agrega el registro
			return nil, fmt.Errorf("usuario no encontrado para tweet %d: %w", tweet.ID, err)
		}
		tweetsWithUser = append(tweetsWithUser, domain.TweetWithUser{
			Tweet:    tweet,
			Username: user.Username,
		})
	}

	return tweetsWithUser, nil
}

// Obtener un tweet por ID
func (repo *TweetRepository) GetTweetByID(tweetID uint) (*domain.Tweet, error) {
	var tweet domain.Tweet
	if err := repo.tweetDB.First(&tweet, tweetID).Error; err != nil {
		return nil, err
	}
	return &tweet, nil
}

// Eliminar un tweet por ID
func (repo *TweetRepository) DeleteTweetByID(tweetID uint) error {
	if err := repo.tweetDB.Delete(&domain.Tweet{}, tweetID).Error; err != nil {
		return errors.New("no se pudo eliminar el tweet")
	}
	return nil
}
