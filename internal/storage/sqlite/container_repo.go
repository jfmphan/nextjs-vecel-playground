package sqlite

import (
	"context"
	"database/sql"

	"homeinventory/internal/domain"
)

// ContainerRepo is the SQLite/libSQL implementation of domain.ContainerRepository.
type ContainerRepo struct{ db *sql.DB }

func NewContainerRepo(db *sql.DB) *ContainerRepo { return &ContainerRepo{db: db} }

func (r *ContainerRepo) Create(ctx context.Context, c *domain.Container) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO containers (name, type, parent_id) VALUES (?, ?, ?)`,
		c.Name, string(c.Type), c.ParentID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *ContainerRepo) Update(ctx context.Context, c *domain.Container) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE containers SET name = ?, type = ?, parent_id = ? WHERE id = ?`,
		c.Name, string(c.Type), c.ParentID, c.ID)
	if err != nil {
		return err
	}
	if affected, err := res.RowsAffected(); err != nil {
		return err
	} else if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Delete removes a container, detaching (rather than deleting) anything that
// referenced it: child containers and items become top-level / unassigned.
func (r *ContainerRepo) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `UPDATE containers SET parent_id = NULL WHERE parent_id = ?`, id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE items SET container_id = NULL WHERE container_id = ?`, id); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `DELETE FROM containers WHERE id = ?`, id)
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

func (r *ContainerRepo) GetByID(ctx context.Context, id int64) (*domain.Container, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, type, parent_id, created_at FROM containers WHERE id = ?`, id)
	c, err := scanContainer(row)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ContainerRepo) List(ctx context.Context) ([]domain.Container, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, type, parent_id, created_at FROM containers ORDER BY name COLLATE NOCASE`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var containers []domain.Container
	for rows.Next() {
		c, err := scanContainer(rows)
		if err != nil {
			return nil, err
		}
		containers = append(containers, c)
	}
	return containers, rows.Err()
}

func (r *ContainerRepo) ItemCounts(ctx context.Context) (map[int64]int, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT container_id, COUNT(*) FROM items
		 WHERE container_id IS NOT NULL GROUP BY container_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[int64]int)
	for rows.Next() {
		var id int64
		var count int
		if err := rows.Scan(&id, &count); err != nil {
			return nil, err
		}
		counts[id] = count
	}
	return counts, rows.Err()
}

func scanContainer(s scanner) (domain.Container, error) {
	var c domain.Container
	var parentID sql.NullInt64
	var typ, createdAt string
	if err := s.Scan(&c.ID, &c.Name, &typ, &parentID, &createdAt); err != nil {
		return domain.Container{}, err
	}
	c.Type = domain.ContainerType(typ)
	c.ParentID = ptrInt64(parentID)
	c.CreatedAt = parseTimestamp(createdAt)
	return c, nil
}
