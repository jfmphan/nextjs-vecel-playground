// Package service holds the application use cases that orchestrate the domain
// repositories. Services own business rules (validation, reference checks,
// derived data) and are independent of any transport concern such as HTTP.
package service

import (
	"context"
	"errors"
	"strings"

	"homeinventory/internal/domain"
)

// ItemInput is the transport-agnostic command for creating or updating an item.
// The HTTP layer maps its request DTOs onto this type.
type ItemInput struct {
	Name              string
	Description       string
	CategoryID        *int64
	ContainerID       *int64
	Quantity          int
	Unit              string
	LowStockThreshold *int
	PurchaseDate      *domain.Date
	ExpiryDate        *domain.Date
	PhotoURL          string
	ValueCents        *int64
	Tags              []string
}

// InventoryService implements the item-related use cases. It depends only on
// repository interfaces, so it can be unit-tested with in-memory fakes.
type InventoryService struct {
	items      domain.ItemRepository
	categories domain.CategoryRepository
	containers domain.ContainerRepository
}

func NewInventoryService(
	items domain.ItemRepository,
	categories domain.CategoryRepository,
	containers domain.ContainerRepository,
) *InventoryService {
	return &InventoryService{items: items, categories: categories, containers: containers}
}

func (s *InventoryService) Create(ctx context.Context, in ItemInput) (*domain.Item, error) {
	if err := s.validate(ctx, in); err != nil {
		return nil, err
	}
	item := in.toDomain(0)
	id, err := s.items.Create(ctx, &item)
	if err != nil {
		return nil, err
	}
	return s.items.GetByID(ctx, id)
}

func (s *InventoryService) Update(ctx context.Context, id int64, in ItemInput) (*domain.Item, error) {
	if err := s.validate(ctx, in); err != nil {
		return nil, err
	}
	item := in.toDomain(id)
	if err := s.items.Update(ctx, &item); err != nil {
		return nil, err
	}
	return s.items.GetByID(ctx, id)
}

func (s *InventoryService) Delete(ctx context.Context, id int64) error {
	return s.items.Delete(ctx, id)
}

func (s *InventoryService) Get(ctx context.Context, id int64) (*domain.Item, error) {
	return s.items.GetByID(ctx, id)
}

func (s *InventoryService) List(ctx context.Context, filter domain.ItemFilter) ([]domain.Item, error) {
	return s.items.List(ctx, filter)
}

func (s *InventoryService) validate(ctx context.Context, in ItemInput) error {
	if strings.TrimSpace(in.Name) == "" {
		return domain.Invalid("name is required")
	}
	if in.Quantity < 0 {
		return domain.Invalid("quantity cannot be negative")
	}
	if in.LowStockThreshold != nil && *in.LowStockThreshold < 0 {
		return domain.Invalid("low stock threshold cannot be negative")
	}
	if in.CategoryID != nil {
		ok, err := s.categories.Exists(ctx, *in.CategoryID)
		if err != nil {
			return err
		}
		if !ok {
			return domain.Invalid("category does not exist")
		}
	}
	if in.ContainerID != nil {
		if _, err := s.containers.GetByID(ctx, *in.ContainerID); err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return domain.Invalid("container does not exist")
			}
			return err
		}
	}
	return nil
}

func (in ItemInput) toDomain(id int64) domain.Item {
	return domain.Item{
		ID:                id,
		Name:              strings.TrimSpace(in.Name),
		Description:       in.Description,
		CategoryID:        in.CategoryID,
		ContainerID:       in.ContainerID,
		Quantity:          in.Quantity,
		Unit:              strings.TrimSpace(in.Unit),
		LowStockThreshold: in.LowStockThreshold,
		PurchaseDate:      in.PurchaseDate,
		ExpiryDate:        in.ExpiryDate,
		PhotoURL:          in.PhotoURL,
		ValueCents:        in.ValueCents,
		Tags:              in.Tags,
	}
}
