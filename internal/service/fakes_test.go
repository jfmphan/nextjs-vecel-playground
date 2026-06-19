package service

import (
	"context"

	"homeinventory/internal/domain"
)

// In-memory repository fakes used by the service tests. Because the services
// depend on the domain repository interfaces (not the SQLite implementations),
// they can be exercised without a database.

type fakeItemRepo struct {
	items  map[int64]domain.Item
	nextID int64
}

func newFakeItemRepo() *fakeItemRepo { return &fakeItemRepo{items: map[int64]domain.Item{}} }

func (r *fakeItemRepo) Create(_ context.Context, item *domain.Item) (int64, error) {
	r.nextID++
	item.ID = r.nextID
	r.items[r.nextID] = *item
	return r.nextID, nil
}

func (r *fakeItemRepo) Update(_ context.Context, item *domain.Item) error {
	if _, ok := r.items[item.ID]; !ok {
		return domain.ErrNotFound
	}
	r.items[item.ID] = *item
	return nil
}

func (r *fakeItemRepo) Delete(_ context.Context, id int64) error {
	if _, ok := r.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.items, id)
	return nil
}

func (r *fakeItemRepo) GetByID(_ context.Context, id int64) (*domain.Item, error) {
	it, ok := r.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &it, nil
}

func (r *fakeItemRepo) List(_ context.Context, _ domain.ItemFilter) ([]domain.Item, error) {
	out := make([]domain.Item, 0, len(r.items))
	for _, it := range r.items {
		out = append(out, it)
	}
	return out, nil
}

type fakeCategoryRepo struct{ existing map[int64]bool }

func (r *fakeCategoryRepo) Create(context.Context, string) (domain.Category, error) {
	return domain.Category{}, nil
}
func (r *fakeCategoryRepo) Delete(context.Context, int64) error          { return nil }
func (r *fakeCategoryRepo) List(context.Context) ([]domain.Category, error) { return nil, nil }
func (r *fakeCategoryRepo) Exists(_ context.Context, id int64) (bool, error) {
	return r.existing[id], nil
}

type fakeContainerRepo struct {
	containers map[int64]domain.Container
	nextID     int64
}

func newFakeContainerRepo() *fakeContainerRepo {
	return &fakeContainerRepo{containers: map[int64]domain.Container{}}
}

func (r *fakeContainerRepo) Create(_ context.Context, c *domain.Container) (int64, error) {
	r.nextID++
	c.ID = r.nextID
	r.containers[r.nextID] = *c
	return r.nextID, nil
}

func (r *fakeContainerRepo) Update(_ context.Context, c *domain.Container) error {
	if _, ok := r.containers[c.ID]; !ok {
		return domain.ErrNotFound
	}
	r.containers[c.ID] = *c
	return nil
}

func (r *fakeContainerRepo) Delete(_ context.Context, id int64) error {
	delete(r.containers, id)
	return nil
}

func (r *fakeContainerRepo) GetByID(_ context.Context, id int64) (*domain.Container, error) {
	c, ok := r.containers[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &c, nil
}

func (r *fakeContainerRepo) List(context.Context) ([]domain.Container, error) {
	out := make([]domain.Container, 0, len(r.containers))
	for _, c := range r.containers {
		out = append(out, c)
	}
	return out, nil
}

func (r *fakeContainerRepo) ItemCounts(context.Context) (map[int64]int, error) {
	return map[int64]int{}, nil
}

func ptr[T any](v T) *T { return &v }
