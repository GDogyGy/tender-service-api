-- +goose Up
-- +goose StatementBegin
CREATE TYPE organization_type AS ENUM ('IE','LLC','JSC');
CREATE TABLE organization
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    type        organization_type,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO organization
VALUES ('6cd13d02-ae0c-40cc-a7ae-967f923bd8cb','BlackRock', 'BlackRock is finance organization', 'IE'),
       ('b4429652-1b66-4f96-a37a-d089e9dbf7a8', 'WaterBlue', 'WATER BLUE focuses on logistics operations in various countries.', 'JSC');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS organization;
DROP TYPE IF EXISTS organization_type;
-- +goose StatementEnd
