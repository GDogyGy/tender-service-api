-- +goose Up
-- +goose StatementBegin
CREATE
    EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE bid_feedback
(
    id          UUID      DEFAULT uuid_generate_v4(),
    bid_id      UUID NOT NULL,
    description TEXT NOT NULL,
    responsible UUID REFERENCES organization_responsible (id),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bid_feedback;
-- +goose StatementEnd
