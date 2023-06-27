package models

import (
	"math/rand"
	"time"
)

type User struct {
	UserID    uint   `json:"user_id" gorm:"primary_key"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"user_name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
	Shortlies []Shortly
	IsGuest   bool       `json:"is_guest" gorm:"default:false"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%^&*()_+=-"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func CreateGuestUser() (*User, error) {

	// create guest user if not exists using email
	var guestUser User
	result := DB.Where("email = ?", "guest@guest.com").Find(&guestUser)
	if result.Error != nil {
		return nil, result.Error
	}

	if guestUser.Email != "guest@guest.com" {
		guestUser = User{
			// override gorm user id to string 'h123'
			UserID:    12,
			FirstName: RandomString(10),
			LastName:  RandomString(10),
			UserName:  RandomString(10),
			Email:     "guest@guest.com",
			Password:  RandomString(10),
			IsGuest:   true,
		}
		result := DB.Create(&guestUser)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		guestUser.Password = RandomString(10)
		result := DB.Save(&guestUser)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return &guestUser, nil
}

// Get guest user
func GetGuestUser() (*User, error) {
	var guestUser User
	result := DB.Where("email = ?", "guest@guest.com").Find(&guestUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &guestUser, nil
}
