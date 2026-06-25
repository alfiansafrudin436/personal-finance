GOPATH := $(shell go env GOPATH)
SQLC_VERSION = v1.27.0

.PHONY: check-sqlc install-sqlc

check-sqlc:
	@if ! command -v sqlc >/dev/null 2>&1; then \
		echo "sqlc not found. Installing sqlc $(SQLC_VERSION)..."; \
		make install-sqlc; \
	else \
		INSTALLED_VERSION=$$(sqlc version | grep -o 'v[0-9]*\.[0-9]*\.[0-9]*'); \
		if [ "$$INSTALLED_VERSION" != "$(SQLC_VERSION)" ]; then \
			echo "sqlc version $$INSTALLED_VERSION found. Updating to $(SQLC_VERSION)..."; \
			make install-sqlc; \
		else \
			echo "sqlc $(SQLC_VERSION) is already installed. Skipping installation."; \
		fi; \
	fi

install-sqlc:
	@echo "Installing sqlc $(SQLC_VERSION)..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)
	@echo "sqlc $(SQLC_VERSION) installed successfully."

.PHONY: init
init:
	make check-sqlc

.PHONY: sqlc
sqlc:
	@sqlc generate

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: update
update:
	@go get -u ./...
	@go mod tidy

.PHONY: build
build:
	@go build -o bin/main .

.PHONY: run
run:
	@go run main.go

.PHONY: test
test:
	@go test -v ./... -cover

.PHONY: docker-build
docker-build:
	@docker build -t personal-finance:latest .

.PHONY: docker-run
docker-run:
	@docker run --env-file .env -p 9000:9000 personal-finance:latest
