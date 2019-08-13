# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)

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

