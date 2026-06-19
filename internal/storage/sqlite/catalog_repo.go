package sqlite

import (
	"context"
	"database/sql"

	"homeinventory/internal/domain"
)

// CategoryRepo is the SQLite/libSQL implementation of domain.CategoryRepository.
type CategoryRepo struct{ db *sql.DB }

func NewCategoryRepo(db *sql.DB) *CategoryRepo { return &CategoryRepo{db: db} }

// Create inserts a category, returning the existing one if the name is already
// taken (names are unique), which keeps the operation idempotent.
func (r *CategoryRepo) Create(ctx context.Context, name string) (domain.Category, error) {
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO categories (name) VALUES (?) ON CONFLICT(name) DO NOTHING`, name); err != nil {
		return domain.Category{}, err
	}
	var c domain.Category
	err := r.db.QueryRowContext(ctx, `SELECT id, name FROM categories WHERE name = ?`, name).
		Scan(&c.ID, &c.Name)
	return c, err
}

// Delete removes a category and unassigns it from any items that used it.
func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `UPDATE items SET category_id = NULL WHERE category_id = ?`, id); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `DELETE FROM categories WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if affected, err := res.RowsAffected(); err != nil {
		return err
	} else if affected == 0 {
		return domain.ErrNotFound
	}
	return tx.Commit()
}

func (r *CategoryRepo) List(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM categories ORDER BY name COLLATE NOCASE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *CategoryRepo) Exists(ctx context.Context, id int64) (bool, error) {
	var one int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM categories WHERE id = ?`, id).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// TagRepo is the SQLite/libSQL implementation of domain.TagRepository.
type TagRepo struct{ db *sql.DB }

func NewTagRepo(db *sql.DB) *TagRepo { return &TagRepo{db: db} }

func (r *TagRepo) List(ctx context.Context) ([]domain.Tag, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM tags ORDER BY name COLLATE NOCASE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}
