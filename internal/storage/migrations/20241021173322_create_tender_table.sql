-- +goose Up
-- +goose StatementBegin
CREATE TYPE tender_status AS ENUM ('CREATED', 'PUBLISHED', 'CLOSED');
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE tender
(
    id          UUID PRIMARY KEY       DEFAULT uuid_generate_v4(),
    name        VARCHAR(100)  NOT NULL,
    description TEXT          NOT NULL,
    service_type VARCHAR(200)  NOT NULL,
    status      tender_status NOT NULL,
    responsible UUID REFERENCES organization_responsible (id)

);

INSERT INTO tender
(name, description, service_type, status, responsible)
VALUES ('Tender: Global Village Construct',
        'Постройка с отделкой комплекса под ключ',
        'Building',
        'CREATED',
        (SELECT organization_responsible.id
         FROM organization_responsible INNER JOIN employee
         ON employee.username = 'user1' AND organization_responsible.user_id = employee.id));

INSERT INTO tender
(name, description, service_type, status, responsible)
VALUES ('Tender: Development Farmers Web App',
        'Разработка приложение для фермеров',
        'Development',
        'PUBLISHED',
        (SELECT organization_responsible.id
         FROM organization_responsible INNER JOIN employee
                                                  ON employee.username = 'user1' AND organization_responsible.user_id = employee.id));

INSERT INTO tender
(name, description, service_type, status, responsible)
VALUES ('Tender: Development Sellers Web App',
        'Разработка приложение для продавцов',
        'Development',
        'PUBLISHED',
        (SELECT organization_responsible.id
         FROM organization_responsible INNER JOIN employee
                                                  ON employee.username = 'user2' AND organization_responsible.user_id = employee.id));

INSERT INTO tender
(name, description, service_type, status, responsible)
VALUES ('Tender: inspect restaurant',
        'Проверить рестораны на качество услуг',
        'Examination',
        'CLOSED',
        (SELECT organization_responsible.id
         FROM organization_responsible INNER JOIN employee
                                                  ON employee.username = 'user2' AND organization_responsible.user_id = employee.id));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tender;
DROP TYPE IF EXISTS tender_status;
-- +goose StatementEnd
