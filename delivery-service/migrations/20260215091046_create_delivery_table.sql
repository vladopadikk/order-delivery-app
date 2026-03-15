-- +goose Up
-- +goose StatementBegin
CREATE TABLE deliveries (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    address TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE deliveries;
-- +goose StatementEnd
