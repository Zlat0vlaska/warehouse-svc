-- +goose Up
CREATE TABLE products (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    price      INTEGER NOT NULL,
    stock      INTEGER NOT NULL
);

-- +goose Down
DROP TABLE products;