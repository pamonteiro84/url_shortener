package storage

import (
	"gorm.io/gorm"
	"url_shortener/internal/models"
)

type URLRepository interface {
	Save(u *models.URL) error
	GetByShortURL(shortURL string) (*models.URL, error)
	Delete(shortCode string) error
}

type gormURLRepository struct {
	db *gorm.DB
}

func NewGormURLRepository(db *gorm.DB) URLRepository {
	return &gormURLRepository{db: db}
}

func (r *gormURLRepository) Save(u *models.URL) error {
	return r.db.Create(u).Error
}

func (r *gormURLRepository) GetByShortURL(shortURL string) (*models.URL, error) {
	var u models.URL
	if err := r.db.Where("short_code = ?", shortURL).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *gormURLRepository) Delete(shortCode string) error {
	result := r.db.Where("short_code = ?", shortCode).Delete(&models.URL{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
