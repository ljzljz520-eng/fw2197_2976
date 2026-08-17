package service

import (
	"errors"
	"strings"

	"skillsassessment/domain"
	"skillsassessment/repository"
)

type CatalogService struct {
	catalog *repository.CatalogRepository
}

func NewCatalogService(catalog *repository.CatalogRepository) *CatalogService {
	return &CatalogService{catalog: catalog}
}

func (s *CatalogService) ListAvailable() ([]repository.CatalogEntry, error) {
	if s == nil || s.catalog == nil {
		return nil, errors.New("catalog is required")
	}
	return s.catalog.Available()
}

func (s *CatalogService) Search(name string) ([]repository.CatalogEntry, error) {
	if s == nil || s.catalog == nil {
		return nil, errors.New("catalog is required")
	}
	return s.catalog.FindByName(strings.TrimSpace(name))
}

func (s *CatalogService) Counts() (map[domain.ProjectStatus]int, error) {
	if s == nil || s.catalog == nil {
		return nil, errors.New("catalog is required")
	}
	return s.catalog.StatusCounts()
}

func IsAvailable(status domain.ProjectStatus) bool {
	return status == domain.ProjectPublished
}
