package service

import (
	"context"

	"homeinventory/internal/domain"
)

// Stats is the dashboard summary of the inventory.
type Stats struct {
	TotalItems    int
	TotalQuantity int
	LowStock      []domain.Item
	Expiring      []domain.Item
}

// StatsService computes the dashboard figures.
type StatsService struct {
	items domain.ItemRepository
}

func NewStatsService(items domain.ItemRepository) *StatsService {
	return &StatsService{items: items}
}

// Compute returns inventory totals plus the items that are low on stock or
// expiring on/before today + expiringWithinDays.
func (s *StatsService) Compute(ctx context.Context, expiringWithinDays int) (Stats, error) {
	all, err := s.items.List(ctx, domain.ItemFilter{})
	if err != nil {
		return Stats{}, err
	}
	cutoff := domain.Today().AddDays(expiringWithinDays)

	stats := Stats{TotalItems: len(all)}
	for _, item := range all {
		stats.TotalQuantity += item.Quantity
		if item.IsLowStock() {
			stats.LowStock = append(stats.LowStock, item)
		}
		if item.ExpiryDate != nil && item.ExpiryDate.Compare(cutoff) <= 0 {
			stats.Expiring = append(stats.Expiring, item)
		}
	}
	return stats, nil
}
