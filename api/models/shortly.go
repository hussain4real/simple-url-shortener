package models

import (
	"net/url"

	"gorm.io/gorm"
)

type Shortly struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primary_key"`
	RedirectURL string `json:"redirect_url" gorm:"not null"`
	ShortURL    string `json:"short_url" gorm:"unique;not null"`
	Visits      uint   `json:"visits"`
	Random      bool   `json:"random"`
	UserID      uint   `json:"user_id" gorm:"null"`
}

// get all shortlies
func GetAllShortlies() ([]Shortly, error) {
	db := DB
	var shortlies []Shortly
	//get all shortlies with their user
	result := db.Find(&shortlies)
	if result.Error != nil {
		return []Shortly{}, result.Error
	}
	return shortlies, nil
}

// get shortly by id
func GetShortlyById(id uint) (Shortly, error) {
	db := DB
	var shortly Shortly
	result := db.Where("id = ?", id).First(&shortly)
	if result.Error != nil {
		return Shortly{}, result.Error
	}
	return shortly, nil
}

// create shortly
func CreateShortly(shortly Shortly) error {
	db := DB
	// create shortly and assign it to the user that created it
	result := db.Create(&shortly)

	return result.Error
}

// update shortly
func UpdateShortly(shortly Shortly) error {
	db := DB
	result := db.Save(&shortly)
	return result.Error
}

// delete shortly
func DeleteShortly(id uint) error {
	db := DB
	result := db.Unscoped().Delete(&Shortly{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// find shortly by short url
func FindByShortlyUrl(url string) (Shortly, error) {
	db := DB
	var shortly Shortly
	result := db.Where("short_url = ?", url).First(&shortly)
	if result.Error != nil {
		return Shortly{}, result.Error
	}
	return shortly, nil
}

// validate shortly is a valid url link
func IsValidURL(RedirectURL string) bool {
	u, err := url.Parse(RedirectURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}
