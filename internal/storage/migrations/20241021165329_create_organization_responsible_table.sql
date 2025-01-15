-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_responsible
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization (id) ON DELETE CASCADE,
    user_id         UUID REFERENCES employee (id) ON DELETE CASCADE
);

INSERT INTO organization_responsible
(organization_id, user_id)
VALUES ((SELECT organization.id
           FROM organization
           WHERE organization.name = 'BlackRock'),
        (SELECT employee.id
         FROM employee
         where employee.username = 'user1' ));

INSERT INTO organization_responsible
(organization_id, user_id)
VALUES ((SELECT organization.id
         FROM organization
         WHERE organization.name = 'WaterBlue'),
        (SELECT employee.id
         FROM employee
         where employee.username = 'user2' ));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS organization_responsible
-- +goose StatementEnd
