-- +goose Up
-- +goose StatementBegin
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE employee
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username   VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name  VARCHAR(50),
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO employee
VALUES ('ab8120b2-faaf-410b-952d-d5e200799d05','user1', 'user_name_1', 'user_last_name_1'),
       ('e846983b-1341-4c61-a2c4-62c2141ccaed', 'user2', 'user_name_2', 'user_last_name_2');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS employee;
-- +goose StatementEnd
