.SILENT:

export TEST_CONTAINER_NAME=db_test
export TEST_APP_NAME=internal-app_test
export TEST_DB_NAME=test

run: linter
	go run ./cmd/tender-service-api/main.go
linter:
	golangci-lint run ./... --config=./.golangci.yaml
testing:
	go test -run TestHandlePing ./internal/handlers/ping/fetch -count=1 && go test ./... -coverprofile cover.out -count=1

test-coverage: testing
	go tool cover -func cover.out | grep total | awk '{print $3}'

create-migrate:
	goose -dir=./internal/storage/migrations postgres "host=${POSTGRES_HOST} user=${POSTGRES_USERNAME} database=${POSTGRES_DATABASE} password=${POSTGRES_PASSWORD}" up

refresh-migrate: reset-migrate
	goose -dir=./internal/storage/migrations postgres "host=${POSTGRES_HOST} user=${POSTGRES_USERNAME} database=${POSTGRES_DATABASE} password=${POSTGRES_PASSWORD}" up

reset-migrate:
	goose -dir=./internal/storage/migrations postgres "host=${POSTGRES_HOST} user=${POSTGRES_USERNAME} database=${POSTGRES_DATABASE} password=${POSTGRES_PASSWORD}" reset

.PHONY: intergration-run
integration-run:
	goose -dir=./internal/storage/migrations postgres "host=localhost user=root port=5434 database=TenderApiTest password=123" up
	sleep 10
	go clean -testcache
	@echo "${BG_GREEN}Run each test integration${RESET}"
	go test -tags=integration -parallel=1 ./test/handlers/create
	go test -tags=integration -parallel=1 ./test/handlers/edit
	go test -tags=integration -parallel=1 ./test/handlers/rollback

swagger-build:
	oapi-codegen -generate types -package transport -o ./internal/handlers/types/transport/transport.go openapi.yaml

