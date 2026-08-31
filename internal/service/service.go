package service

import (
	"errors"
	"url_shortener/internal/models"
	"url_shortener/internal/shortcode"
	"url_shortener/internal/storage"
)

var ErrAlreadyExists = errors.New("url já tem shortcode")

type Service struct {
	storage storage.URLRepository
}

func NewService(storage storage.URLRepository) *Service {
	return &Service{storage: storage}
}

func (s *Service) ShortenURL(originalURL string) (string, error) {
	shortCode := shortcode.Generate(originalURL)

	if existingURL, err := s.GetOriginalURL(shortCode); err == nil &&
		existingURL == originalURL {
		return shortCode, ErrAlreadyExists
	}
	url := &models.URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
	}
	if err := s.storage.Save(url); err != nil {
		return "", err
	}
	return shortCode, nil
}

func (s *Service) GetOriginalURL(shortCode string) (string, error) {
	url, err := s.storage.GetByShortURL(shortCode)
	if err != nil {
		return "", err
	}
	return url.OriginalURL, nil
}