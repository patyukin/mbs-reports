-- +goose Up
CREATE TABLE IF NOT EXISTS transactions
(
    id                  LowCardinality(String),
    payment_id          LowCardinality(String),
    account_id          LowCardinality(String),
    user_id             LowCardinality(String),
    type                LowCardinality(String),
    amount              Int64,
    currency            LowCardinality(String),
    description         String,
    payment_description String,
    status              LowCardinality(String),
    send_status         LowCardinality(String),
    created_at          DateTime
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(created_at)
ORDER BY (created_at, id);

-- +goose Down
DROP TABLE IF EXISTS transactions;
