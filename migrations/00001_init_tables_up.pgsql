-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS balance (
    user_id BIGSERIAL PRIMARY KEY,
    balance DECIMAL(15,2) NOT NULL DEFAULT 0 CHECK (balance >= 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE balance;
-- +goose StatementEnd
