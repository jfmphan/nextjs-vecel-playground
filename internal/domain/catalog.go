package domain

import "context"

// Category is a controlled classification for items (a managed vocabulary,
// created explicitly by the user).
type Category struct {
	ID   int64
	Name string
}

// Tag is a free-form label. Tags are created implicitly when attached to an
// item (see ItemRepository), so there is no explicit "create tag" use case.
type Tag struct {
	ID   int64
	Name string
}

// CategoryRepository persists the managed category vocabulary.
type CategoryRepository interface {
	Create(ctx context.Context, name string) (Category, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]Category, error)
	Exists(ctx context.Context, id int64) (bool, error)
}

// TagRepository reads the set of known tags.
type TagRepository interface {
	List(ctx context.Context) ([]Tag, error)
}
