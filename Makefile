.SILENT:

export TEST_CONTAINER_NAME=db_test
export TEST_DB_NAME=test

# Когда тут ставлю test не работает почему-то
testing:
	go test ./... -coverprofile cover.out

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
	docker run --rm -d --name ${TEST_CONTAINER_NAME} -p 5434:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123 -e POSTGRES_DB=TenderApiTest -d postgres:latest
	sleep 5
	go clean -testcache
	@echo "${BG_GREEN}Run each test integration${RESET}"
	go test -tags=integration -parallel=1 ./test/handlers/create
	go test -tags=integration -parallel=1 ./test/handlers/edit
	go test -tags=integration -parallel=1 ./test/handlers/rollback
	docker stop ${TEST_CONTAINER_NAME}
