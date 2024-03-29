NAME=sort-MP3

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20


define colored
	@echo '${GREEN}$1${RESET}'
endef

## Show help
help:
	${call colored, help is running...}
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## dependencies - fetch all dependencies for scripts
dependencies:
	${call colored, dependensies is running...}
	./scripts/get-dependencies.sh

## lint project
lint:
	${call colored, lint is running...}
	./scripts/linters.sh
.PHONY: lint

## Test all packages
test:
	${call colored, test is running...}
	./scripts/tests.sh
.PHONY: test

## Test coverage
test-cover:
	${call colored, test-cover is running...}
	go test -race -coverpkg=./... -v -coverprofile .testCoverage.out ./...
	gocov convert .testCoverage.out | gocov report
.PHONY: test-cover

new-version: lint test compile
	${call colored, new version is running...}
	./scripts/version.sh
.PHONY: new-version


## Formats the code.
format:
	${call colored, formatting is running...}
	go vet ./...
	go fmt ./...

## Fix-imports order
fix-imports:
	${call colored, fixing imports...}
	./scripts/fix-imports-order.sh

## DB migration commands:
## Create migration files with name which should be specified with flag 'n=some_name'.
create-migrations:
	migrate create -ext sql -dir migrations -seq $(n)
	${call colored,migrations is created}

## "postgres://user:password@host:port/name_db?sslmode=disable"
database=postgres://sorter:master@localhost:5433/finndon?sslmode=disable

## Roll migrations.
migrate-up:
	migrate -path ./migrations -database $(database) up
	${call colored,migrations is upped}

## Rollback all migrations.
## If you specify flag 's=i' this will rollback 'i' migrations.
migrate-down:
	migrate -path ./migrations -database $(database) down $(s)
	${call colored,migrations is downed}

## Drop migrations.
migrate-drop:
	migrate -path ./migrations -database $(database) drop
	${call colored,migrations is droped}
