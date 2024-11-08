package domain

import "time"

type Tweet struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Content   string    `gorm:"size:280"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
}

// Estructura para enriquecer un tweet con datos de usuario
type TweetWithUser struct {
	Tweet
	Username string
}
