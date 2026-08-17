-- +goose Up
CREATE TABLE IF NOT EXISTS messages_stream.sagas
(
	id text NOT NULL,
	name text NOT NULL,
	data bytea NOT NULL,
	step int NOT NULL,
	done bool NOT NULL,
	compensating bool NOT NULL,
	PRIMARY KEY (id, name)
);

-- +goose Down
DROP TABLE IF EXISTS messages_stream.sagas;