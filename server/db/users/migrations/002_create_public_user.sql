-- +goose Up
INSERT INTO users.users(id, login, salt) VALUES (9999999, 'public', '');

-- +goose Down
DELETE FROM users.users WHERE id = 9999999;
