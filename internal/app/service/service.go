package service

import (
	"fmt"
	"log/slog"
)

type Repository interface {
	SaveUrl(string, string) (int64, error)
	GetUrl(string) (string, error)
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) SaveUrl(urlToSave string, alias string) error {
	slog.Debug(fmt.Sprintf("saving url: %s with alias: %s", urlToSave, alias))
	id, err := s.repository.SaveUrl(urlToSave, alias)
	if err != nil {
		return fmt.Errorf("couldn't save url: %w", err)
	}
	slog.Debug(fmt.Sprintf("saved with id: %d", id))
	return nil
}

func (s *Service) GetUrl(alias string) (string, error) {
	slog.Debug(fmt.Sprintf("Get url for alias: %s", alias))
	resultUrl, err := s.repository.GetUrl(alias)
	if err != nil {
		return "", fmt.Errorf("couldn't get url: %w", err)
	}
	slog.Debug(fmt.Sprintf("got url: %s for alias: %s", resultUrl, alias))
	return resultUrl, nil
}
