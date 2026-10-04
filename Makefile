COVERAGE_MIN ?= 80
COVERAGE_FILE := coverage.out
# Пакет main (cmd/gendiff) из-под порога исключаем: в нём только вызов app.Run.
COVERAGE_PKGS := $(shell go list ./... | grep -v /cmd/)

# Версия линтера закреплена здесь и только здесь: и локальный запуск, и CI
# идут через `make lint`, поэтому один коммит получает один и тот же результат.
# Обновление линтера — отдельное изменение этой строки.
GOLANGCI_LINT_VERSION := v2.12.1
# Версия в имени файла: после её смены make поставит новый бинарник сам.
GOLANGCI_LINT := bin/golangci-lint-$(GOLANGCI_LINT_VERSION)

build:
	go build -o bin/gendiff ./cmd/gendiff

$(GOLANGCI_LINT):
	GOBIN=$(CURDIR)/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	mv bin/golangci-lint $(GOLANGCI_LINT)

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

lint-fix: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --fix ./...

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
