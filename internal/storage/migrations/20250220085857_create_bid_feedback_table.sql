-- +goose Up
-- +goose StatementBegin
CREATE
    EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE bid_feedback
(
    id          UUID      DEFAULT uuid_generate_v4(),
--     TODO: получается bid_id без связи спросить у димы ок ли это?
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
