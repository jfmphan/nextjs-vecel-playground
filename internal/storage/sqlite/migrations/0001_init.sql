CREATE TABLE categories (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE containers (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    type       TEXT NOT NULL,
    parent_id  INTEGER REFERENCES containers(id),
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE items (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    name                TEXT NOT NULL,
    description         TEXT NOT NULL DEFAULT '',
    category_id         INTEGER REFERENCES categories(id),
    container_id        INTEGER REFERENCES containers(id),
    quantity            INTEGER NOT NULL DEFAULT 1,
    unit                TEXT NOT NULL DEFAULT '',
    low_stock_threshold INTEGER,
    purchase_date       TEXT,
    expiry_date         TEXT,
    photo_url           TEXT NOT NULL DEFAULT '',
    value_cents         INTEGER,
    created_at          TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at          TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE tags (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE item_tags (
    item_id INTEGER NOT NULL REFERENCES items(id),
    tag_id  INTEGER NOT NULL REFERENCES tags(id),
    PRIMARY KEY (item_id, tag_id)
);

CREATE INDEX idx_items_container ON items(container_id);
CREATE INDEX idx_items_category ON items(category_id);
CREATE INDEX idx_items_expiry ON items(expiry_date);
CREATE INDEX idx_containers_parent ON containers(parent_id);
