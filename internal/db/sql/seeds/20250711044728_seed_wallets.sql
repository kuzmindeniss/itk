-- +goose Up
INSERT INTO wallets (id, balance) VALUES ('8e3449a8-5cbc-4159-a8e2-45eea1eebdb1', 0);
INSERT INTO wallets (id, balance) VALUES ('8e3449a8-5cbc-4159-a8e2-45eea1eebdb2', 0);

-- +goose Down
DELETE FROM wallets WHERE id = '8e3449a8-5cbc-4159-a8e2-45eea1eebdb1';
DELETE FROM wallets WHERE id = '8e3449a8-5cbc-4159-a8e2-45eea1eebdb2';

