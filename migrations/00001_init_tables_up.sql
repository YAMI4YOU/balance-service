-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wallet (
    user_id BIGINT PRIMARY KEY,
    balance NUMERIC(15, 2) NOT NULL DEFAULT 0 CHECK(balance >= 0)
);

CREATE TABLE IF NOT EXISTS reservation (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    service_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'reserved' CHECK (status IN ('reserved', 'confirmed', 'canceled'))
);

CREATE INDEX IF NOT EXISTS idx_reservation_user_id ON reservation(user_id);

CREATE INDEX IF NOT EXISTS idx_reservation_order_id ON reservation(order_id);

CREATE INDEX IF NOT EXISTS idx_reservation_service_id ON reservation(service_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reservation;
DROP TABLE wallet;
-- +goose StatementEnd
