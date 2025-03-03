-- +goose Up
-- +goose StatementBegin
CREATE TYPE decision_status AS ENUM ('REJECTED', 'APPROVED');
CREATE
    EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE decision_bid
(
    id          UUID      DEFAULT uuid_generate_v4(),
    bid_id      UUID NOT NULL,
    decision    decision_status NOT NULL,
    responsible UUID REFERENCES organization_responsible (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS decision_bid;
DROP TYPE IF EXISTS decision_status;
-- +goose StatementEnd
