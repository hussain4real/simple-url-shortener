package models

import (
	"gorm.io/gorm"
)

type Shortly struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primary_key"`
	RedirectURL string `json:"redirect_url" gorm:"not null"`
	ShortURL    string `json:"short_url" gorm:"unique;not null"`
	Visits      uint   `json:"visits"`
	Random      bool   `json:"random"`
}

// get all shortlies
func GetAllShortlies() ([]Shortly, error) {
	db := DB
	var shortly []Shortly
	result := db.Find(&shortly)
	if result.Error != nil {
		return []Shortly{}, result.Error
	}
	return shortly, nil
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
func CreateShortly(shortly Shortly) (Shortly, error) {
	db := DB
	result := db.Create(&shortly)
	if result.Error != nil {
		return Shortly{}, result.Error
	}
	return shortly, nil
}

// update shortly
func UpdateShortly(shortly Shortly) (Shortly, error) {
	db := DB
	result := db.Save(&shortly)
	if result.Error != nil {
		return Shortly{}, result.Error
	}
	return shortly, nil
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
func findByShortlyUrl(url string) (Shortly, error) {
	db := DB
	var shortly Shortly
	result := db.Where("short_url = ?", url).First(&shortly)
	if result.Error != nil {
		return Shortly{}, result.Error
	}
	return shortly, nil
}
