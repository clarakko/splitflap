package service

import (
	"splitflap-api-go/internal/model"
	"splitflap-api-go/internal/repository"
)

type DisplayService struct {
	repo repository.DisplayRepository
}

func NewDisplayService(repo repository.DisplayRepository) *DisplayService {
	return &DisplayService{
		repo: repo,
	}
}

func (s *DisplayService) GetDisplay(id string) *model.Display {
	display, err := s.repo.GetByID(id)
	if err != nil {
		return nil
	}
	return display
}

func (s *DisplayService) CreateDisplay(display *model.Display) error {
	return s.repo.Create(display)
}

func (s *DisplayService) UpdateDisplay(display *model.Display) error {
	return s.repo.Update(display)
}

func (s *DisplayService) DeleteDisplay(id string) error {
	return s.repo.Delete(id)
}

func (s *DisplayService) ListDisplays() ([]*model.Display, error) {
	return s.repo.List()
}
