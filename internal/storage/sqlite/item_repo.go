package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"homeinventory/internal/domain"
)

// ItemRepo is the SQLite/libSQL implementation of domain.ItemRepository.
// It persists an item together with its tags as a single aggregate.
type ItemRepo struct{ db *sql.DB }

func NewItemRepo(db *sql.DB) *ItemRepo { return &ItemRepo{db: db} }

const itemColumns = `id, name, description, category_id, container_id, quantity, unit,
	low_stock_threshold, purchase_date, expiry_date, photo_url, value_cents,
	created_at, updated_at`

func (r *ItemRepo) Create(ctx context.Context, item *domain.Item) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `INSERT INTO items
		(name, description, category_id, container_id, quantity, unit,
		 low_stock_threshold, purchase_date, expiry_date, photo_url, value_cents)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.Name, item.Description, item.CategoryID, item.ContainerID, item.Quantity,
		item.Unit, item.LowStockThreshold, dateArg(item.PurchaseDate), dateArg(item.ExpiryDate),
		item.PhotoURL, item.ValueCents)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if err := replaceTags(ctx, tx, id, item.Tags); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ItemRepo) Update(ctx context.Context, item *domain.Item) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `UPDATE items SET
		name = ?, description = ?, category_id = ?, container_id = ?, quantity = ?,
		unit = ?, low_stock_threshold = ?, purchase_date = ?, expiry_date = ?,
		photo_url = ?, value_cents = ?, updated_at = datetime('now')
		WHERE id = ?`,
		item.Name, item.Description, item.CategoryID, item.ContainerID, item.Quantity,
		item.Unit, item.LowStockThreshold, dateArg(item.PurchaseDate), dateArg(item.ExpiryDate),
		item.PhotoURL, item.ValueCents, item.ID)
	if err != nil {
		return err
	}
	if affected, err := res.RowsAffected(); err != nil {
		return err
	} else if affected == 0 {
		return domain.ErrNotFound
	}
	if err := replaceTags(ctx, tx, item.ID, item.Tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *ItemRepo) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM item_tags WHERE item_id = ?`, id); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `DELETE FROM items WHERE id = ?`, id)
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

func (r *ItemRepo) GetByID(ctx context.Context, id int64) (*domain.Item, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+itemColumns+` FROM items WHERE id = ?`, id)
	item, err := scanItem(row)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	tagsByItem, err := r.tagsFor(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	item.Tags = tagsByItem[id]
	return &item, nil
}

func (r *ItemRepo) List(ctx context.Context, filter domain.ItemFilter) ([]domain.Item, error) {
	where, args := buildItemFilter(filter)
	query := `SELECT ` + itemColumns + ` FROM items`
	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, ` AND `)
	}
	query += ` ORDER BY name COLLATE NOCASE`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := r.attachTags(ctx, items); err != nil {
		return nil, err
	}
	return items, nil
}

// buildItemFilter turns an ItemFilter into SQL WHERE fragments and their args.
func buildItemFilter(f domain.ItemFilter) (where []string, args []any) {
	if f.Query != "" {
		where = append(where, `(LOWER(name) LIKE ? OR LOWER(description) LIKE ?)`)
		pattern := "%" + strings.ToLower(f.Query) + "%"
		args = append(args, pattern, pattern)
	}
	if f.CategoryID != nil {
		where = append(where, `category_id = ?`)
		args = append(args, *f.CategoryID)
	}
	if f.ContainerID != nil {
		where = append(where, `container_id = ?`)
		args = append(args, *f.ContainerID)
	}
	if f.LowStockOnly {
		where = append(where, `low_stock_threshold IS NOT NULL AND quantity <= low_stock_threshold`)
	}
	if f.ExpiringBefore != nil {
		where = append(where, `expiry_date IS NOT NULL AND expiry_date <= ?`)
		args = append(args, f.ExpiringBefore.String())
	}
	if f.Tag != "" {
		where = append(where, `id IN (SELECT it.item_id FROM item_tags it
			JOIN tags t ON t.id = it.tag_id WHERE t.name = ?)`)
		args = append(args, f.Tag)
	}
	return where, args
}

func scanItem(s scanner) (domain.Item, error) {
	var (
		item                                          domain.Item
		categoryID, containerID, lowStock, valueCents sql.NullInt64
		purchase, expiry                              domain.Date
		createdAt, updatedAt                          string
	)
	err := s.Scan(&item.ID, &item.Name, &item.Description, &categoryID, &containerID,
		&item.Quantity, &item.Unit, &lowStock, &purchase, &expiry,
		&item.PhotoURL, &valueCents, &createdAt, &updatedAt)
	if err != nil {
		return domain.Item{}, err
	}
	item.CategoryID = ptrInt64(categoryID)
	item.ContainerID = ptrInt64(containerID)
	item.LowStockThreshold = ptrInt(lowStock)
	item.ValueCents = ptrInt64(valueCents)
	item.PurchaseDate = optionalDate(purchase)
	item.ExpiryDate = optionalDate(expiry)
	item.CreatedAt = parseTimestamp(createdAt)
	item.UpdatedAt = parseTimestamp(updatedAt)
	return item, nil
}

func (r *ItemRepo) attachTags(ctx context.Context, items []domain.Item) error {
	if len(items) == 0 {
		return nil
	}
	ids := make([]int64, len(items))
	for i := range items {
		ids[i] = items[i].ID
	}
	byItem, err := r.tagsFor(ctx, ids)
	if err != nil {
		return err
	}
	for i := range items {
		items[i].Tags = byItem[items[i].ID]
	}
	return nil
}

// tagsFor loads the tag names for the given item IDs in a single query,
// avoiding an N+1 lookup when listing items.
func (r *ItemRepo) tagsFor(ctx context.Context, ids []int64) (map[int64][]string, error) {
	byItem := make(map[int64][]string)
	if len(ids) == 0 {
		return byItem, nil
	}
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := `SELECT it.item_id, t.name FROM item_tags it
		JOIN tags t ON t.id = it.tag_id
		WHERE it.item_id IN (` + strings.Join(placeholders, ",") + `)
		ORDER BY t.name COLLATE NOCASE`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var itemID int64
		var name string
		if err := rows.Scan(&itemID, &name); err != nil {
			return nil, err
		}
		byItem[itemID] = append(byItem[itemID], name)
	}
	return byItem, rows.Err()
}

// replaceTags sets an item's tags to exactly the given names, creating any tag
// that does not yet exist. It runs inside the caller's transaction so the item
// and its tags are persisted atomically.
func replaceTags(ctx context.Context, tx *sql.Tx, itemID int64, names []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM item_tags WHERE item_id = ?`, itemID); err != nil {
		return err
	}
	seen := make(map[string]bool)
	for _, raw := range names {
		name := strings.TrimSpace(raw)
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true

		if _, err := tx.ExecContext(ctx,
			`INSERT INTO tags (name) VALUES (?) ON CONFLICT(name) DO NOTHING`, name); err != nil {
			return err
		}
		var tagID int64
		if err := tx.QueryRowContext(ctx, `SELECT id FROM tags WHERE name = ?`, name).Scan(&tagID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO item_tags (item_id, tag_id) VALUES (?, ?)`, itemID, tagID); err != nil {
			return err
		}
	}
	return nil
}
