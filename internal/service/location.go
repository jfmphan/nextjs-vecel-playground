package service

import (
	"context"
	"errors"
	"strings"

	"homeinventory/internal/domain"
)

// ContainerInput is the transport-agnostic command for creating or updating a
// container (storage location).
type ContainerInput struct {
	Name     string
	Type     domain.ContainerType
	ParentID *int64
}

// LocationService implements the container/location use cases, including
// nesting rules.
type LocationService struct {
	containers domain.ContainerRepository
}

func NewLocationService(containers domain.ContainerRepository) *LocationService {
	return &LocationService{containers: containers}
}

func (s *LocationService) Create(ctx context.Context, in ContainerInput) (*domain.Container, error) {
	in.Type = normalizeType(in.Type)
	if err := s.validate(ctx, 0, in); err != nil {
		return nil, err
	}
	c := domain.Container{Name: strings.TrimSpace(in.Name), Type: in.Type, ParentID: in.ParentID}
	id, err := s.containers.Create(ctx, &c)
	if err != nil {
		return nil, err
	}
	return s.containers.GetByID(ctx, id)
}

func (s *LocationService) Update(ctx context.Context, id int64, in ContainerInput) (*domain.Container, error) {
	in.Type = normalizeType(in.Type)
	if err := s.validate(ctx, id, in); err != nil {
		return nil, err
	}
	c := domain.Container{ID: id, Name: strings.TrimSpace(in.Name), Type: in.Type, ParentID: in.ParentID}
	if err := s.containers.Update(ctx, &c); err != nil {
		return nil, err
	}
	return s.containers.GetByID(ctx, id)
}

func (s *LocationService) Delete(ctx context.Context, id int64) error {
	return s.containers.Delete(ctx, id)
}

func (s *LocationService) Get(ctx context.Context, id int64) (*domain.Container, error) {
	return s.containers.GetByID(ctx, id)
}

func (s *LocationService) List(ctx context.Context) ([]domain.Container, error) {
	return s.containers.List(ctx)
}

// ItemCounts reports how many items sit directly in each container, keyed by
// container ID.
func (s *LocationService) ItemCounts(ctx context.Context) (map[int64]int, error) {
	return s.containers.ItemCounts(ctx)
}

func (s *LocationService) validate(ctx context.Context, id int64, in ContainerInput) error {
	if strings.TrimSpace(in.Name) == "" {
		return domain.Invalid("name is required")
	}
	if !in.Type.Valid() {
		return domain.Invalid("invalid container type")
	}
	return s.ensureNoCycle(ctx, id, in.ParentID)
}

// ensureNoCycle walks the proposed parent chain to confirm the parent exists and
// that re-parenting would not create a loop (a container cannot be nested under
// itself or any of its own descendants). For Create, id is 0 and never matches a
// real container, so the walk only validates that the parent exists.
func (s *LocationService) ensureNoCycle(ctx context.Context, id int64, parentID *int64) error {
	if parentID == nil {
		return nil
	}
	if *parentID == id {
		return domain.Invalid("a container cannot be its own parent")
	}
	cursor := parentID
	for cursor != nil {
		parent, err := s.containers.GetByID(ctx, *cursor)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return domain.Invalid("parent container does not exist")
			}
			return err
		}
		if parent.ID == id {
			return domain.Invalid("container hierarchy would create a cycle")
		}
		cursor = parent.ParentID
	}
	return nil
}

func normalizeType(t domain.ContainerType) domain.ContainerType {
	if t == "" {
		return domain.ContainerOther
	}
	return t
}
