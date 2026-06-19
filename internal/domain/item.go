package domain

import (
	"context"
	"time"
)

// Item is the central aggregate: a thing tracked in the inventory. Its tags are
// part of the aggregate and are persisted together with the item.
type Item struct {
	ID                int64
	Name              string
	Description       string
	CategoryID        *int64
	ContainerID       *int64
	Quantity          int
	Unit              string
	LowStockThreshold *int
	PurchaseDate      *Date
	ExpiryDate        *Date
	PhotoURL          string
	ValueCents        *int64
	Tags              []string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// IsLowStock reports whether the item is at or below its low-stock threshold.
// Items without a threshold are never considered low stock.
func (i Item) IsLowStock() bool {
	return i.LowStockThreshold != nil && i.Quantity <= *i.LowStockThreshold
}

// ItemFilter describes optional criteria for listing items. Zero-valued fields
// are ignored, so an empty filter lists everything.
type ItemFilter struct {
	Query          string // matches name/description, case-insensitive
	CategoryID     *int64
	ContainerID    *int64
	Tag            string
	LowStockOnly   bool
	ExpiringBefore *Date
}

// ItemRepository persists Item aggregates, including their tags.
type ItemRepository interface {
	Create(ctx context.Context, item *Item) (int64, error)
	Update(ctx context.Context, item *Item) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Item, error)
	List(ctx context.Context, filter ItemFilter) ([]Item, error)
}
