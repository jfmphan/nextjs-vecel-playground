package domain

import (
	"context"
	"time"
)

// ContainerType classifies a storage location.
type ContainerType string

const (
	ContainerRoom  ContainerType = "room"
	ContainerShelf ContainerType = "shelf"
	ContainerBox   ContainerType = "box"
	ContainerOther ContainerType = "other"
)

// Valid reports whether the type is one of the known values.
func (t ContainerType) Valid() bool {
	switch t {
	case ContainerRoom, ContainerShelf, ContainerBox, ContainerOther:
		return true
	default:
		return false
	}
}

// Container is an optionally nested place where items are stored, e.g.
// Garage → Shelf B → Box 3. Nesting is expressed via ParentID (nil = top level).
type Container struct {
	ID        int64
	Name      string
	Type      ContainerType
	ParentID  *int64
	CreatedAt time.Time
}

// ContainerRepository persists containers and answers tree-related queries.
type ContainerRepository interface {
	Create(ctx context.Context, c *Container) (int64, error)
	Update(ctx context.Context, c *Container) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*Container, error)
	List(ctx context.Context) ([]Container, error)
	// ItemCounts returns the number of items directly in each container, keyed by
	// container ID. Containers with no items are omitted.
	ItemCounts(ctx context.Context) (map[int64]int, error)
}
