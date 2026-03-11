-- +goose Up
CREATE TABLE IF NOT EXISTS users.premiums
(
	id            bigint       NOT NULL,
	invoice_id    VARCHAR(256) NOT NULL,
	created_at    timestamptz  NOT NULL DEFAULT NOW(),
	expires_at    timestamptz  NOT NULL,
	PRIMARY KEY(id, invoice_id)
);

-- +goose Down
DROP TABLE IF EXISTS users.premiums;
