COVERAGE_MIN ?= 80
COVERAGE_FILE := coverage.out
# Пакет main (cmd/gendiff) из-под порога исключаем: в нём только вызов app.Run.
COVERAGE_PKGS := $(shell go list ./... | grep -v /cmd/)

build:
	go build -o bin/gendiff ./cmd/gendiff

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=$(COVERAGE_FILE) $(COVERAGE_PKGS)
	@total=$$(go tool cover -func=$(COVERAGE_FILE) | awk '/^total:/ {sub("%", "", $$3); print $$3}'); \
	echo "Total coverage: $$total% (min $(COVERAGE_MIN)%)"; \
	if awk -v t="$$total" -v m="$(COVERAGE_MIN)" 'BEGIN {exit !(t < m)}'; then \
		echo "Coverage is below threshold"; \
		exit 1; \
	fi

.PHONY: build lint lint-fix test test-coverage
