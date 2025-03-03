-- +goose Up
-- +goose StatementBegin
CREATE TYPE bids_status AS ENUM ('CREATED', 'PUBLISHED', 'CANCELED', 'APPROVED');
CREATE
    EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE bids
(
    id          UUID                  DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    description TEXT         NOT NULL,
    status      bids_status  NOT NULL,
    tender_id   UUID         NOT NULL,
    version     INT          NOT NULL default 1,
    responsible UUID REFERENCES organization_responsible (id)
);

INSERT INTO bids
    (name, description, status, tender_id, version, responsible)
VALUES ('Предложение 1',
        'Построю с отделкой комплекс под ключ',
        'CREATED',
        (SELECT tender.id
         FROM tender
                  INNER JOIN organization_responsible
                             ON organization_responsible.organization_id = tender.responsible
                  INNER JOIN employee
                             ON employee.id = organization_responsible.user_id AND employee.username = 'user2'
         LIMIT 1),
        1,
        (SELECT organization_responsible.id
         FROM organization_responsible
                  INNER JOIN employee
                             ON employee.username = 'user2' AND organization_responsible.user_id = employee.id));
INSERT INTO bids
    (name, description, status, tender_id, version, responsible)
VALUES ('Предложение 2',
        'Разработаю приложение для фермеров',
        'CREATED',
        (SELECT tender.id
         FROM tender
                  INNER JOIN organization_responsible
                             ON organization_responsible.organization_id = tender.responsible
                  INNER JOIN employee
                             ON employee.id = organization_responsible.user_id AND employee.username = 'user1'
         LIMIT 1),
        1,
        (SELECT organization_responsible.id
         FROM organization_responsible
                  INNER JOIN employee
                             ON employee.username = 'user1' AND organization_responsible.user_id = employee.id));
INSERT INTO bids
    (name, description, status, tender_id, version, responsible)
VALUES ('Предложение 3',
        'Проверю рестораны на качество услуг',
        'CREATED',
        (SELECT tender.id
         FROM tender
                  INNER JOIN organization_responsible
                             ON organization_responsible.organization_id = tender.responsible
                  INNER JOIN employee
                             ON employee.id = organization_responsible.user_id AND employee.username = 'user1'
         LIMIT 1),
        1,
        (SELECT organization_responsible.id
         FROM organization_responsible
                  INNER JOIN employee
                             ON employee.username = 'user1' AND organization_responsible.user_id = employee.id));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bids;
DROP TYPE IF EXISTS bids_status;
-- +goose StatementEnd
