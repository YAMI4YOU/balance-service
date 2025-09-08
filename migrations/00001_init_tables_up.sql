-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wallet (
    user_id BIGINT PRIMARY KEY,
    balance BIGINT NOT NULL DEFAULT 0 CHECK(balance >= 0)
);

DO $$ BEGIN
    CREATE TYPE reservation_status AS ENUM (
        'reserved',
        'confirmed',
        'cancelled'
        );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS reservation (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    service_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    status reservation_status NOT NULL DEFAULT 'reserved'
);

CREATE TABLE IF NOT EXISTS accounting_report (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    service_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    operation_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reservation_user_id ON reservation(user_id);

CREATE INDEX IF NOT EXISTS idx_reservation_order_id ON reservation(order_id);

CREATE INDEX IF NOT EXISTS idx_reservation_service_id ON reservation(service_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_reservation_service_order_unique
    ON reservation (service_id, order_id);

CREATE INDEX IF NOT EXISTS idx_accounting_report_created_at ON accounting_report(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accounting_report;
DROP TABLE reservation;
DROP TABLE wallet;
-- +goose StatementEnd
