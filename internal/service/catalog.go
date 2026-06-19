package service

import (
	"context"
	"strings"

	"homeinventory/internal/domain"
)

// CatalogService implements the use cases for the supporting vocabularies:
// categories (managed) and tags (read-only here; created via items).
type CatalogService struct {
	categories domain.CategoryRepository
	tags       domain.TagRepository
}

func NewCatalogService(categories domain.CategoryRepository, tags domain.TagRepository) *CatalogService {
	return &CatalogService{categories: categories, tags: tags}
}

func (s *CatalogService) CreateCategory(ctx context.Context, name string) (domain.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Category{}, domain.Invalid("name is required")
	}
	return s.categories.Create(ctx, name)
}

func (s *CatalogService) DeleteCategory(ctx context.Context, id int64) error {
	return s.categories.Delete(ctx, id)
}

func (s *CatalogService) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return s.categories.List(ctx)
}

func (s *CatalogService) ListTags(ctx context.Context) ([]domain.Tag, error) {
	return s.tags.List(ctx)
}
