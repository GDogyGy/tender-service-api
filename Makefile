.SILENT:

export TEST_CONTAINER_NAME=db_test
export TEST_APP_NAME=internal-app_test
export TEST_DB_NAME=test

run: linter
	go run ./cmd/tender-service-api/main.go
linter:
	golangci-lint run ./... --config=./.golangci.yaml
testing:
	go test -run TestHandlePing ./internal/handlers/rest/ping/fetch -count=1 && go test ./... -coverprofile cover.out -count=1

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
	# TODO: получается бинарник тестовой приложухи надо пересобирать через docker-compose чтобы тестить актуальное приложение -_-? и по другому не решить это
	go clean -testcache
	@echo "${BG_GREEN}Run each test integration${RESET}"
	go test -tags=integration -parallel=1 ./integration_tests/handlers/ping/fetch
	go test -tags=integration -parallel=1 ./integration_tests/handlers/tender/create
	go test -tags=integration -parallel=1 ./integration_tests/handlers/tender/fetch
	go test -tags=integration -parallel=1 ./integration_tests/handlers/tender/update
	go test -tags=integration -parallel=1 ./integration_tests/handlers/bids/create
	go test -tags=integration -parallel=1 ./integration_tests/handlers/bids/fetch
	go test -tags=integration -parallel=1 ./integration_tests/handlers/bids/update

swagger-build:
	oapi-codegen -generate types -package transport -o ./internal/handlers/rest/types/transport/transport.go openapi.yaml

proto-gen:
	protoc --go_out=./internal/protos --go-grpc_out=./internal/protos ./internal/protos/proto/tender/proto_tender_service_create.proto
	protoc --go_out=./internal/protos --go-grpc_out=./internal/protos ./internal/protos/proto/tender/proto_tender_service_fetch.proto
	protoc --go_out=./internal/protos --go-grpc_out=./internal/protos ./internal/protos/proto/tender/proto_tender_service_update.proto
	protoc --go_out=./internal/protos --go-grpc_out=./internal/protos ./internal/protos/proto/bids/proto_bids_service_create.proto
	protoc --go_out=./internal/protos --go-grpc_out=./internal/protos ./internal/protos/proto/bids/proto_bids_service_fetch.proto
	protoc --go_out=./internal/protos --go-grpc_out=./internal/protos ./internal/protos/proto/bids/proto_bids_service_update.proto
