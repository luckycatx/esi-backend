APP_NAME = esi-server

all: update gen build
	@echo "Run 'make init' for the first time"
.PHONY: init
init: update sqlc wire

.PHONY: build
build:
	go build -o ./bin/${APP_NAME}.exe -trimpath ./cmd/...

.PHONY: build-all
build-all:
	go build -o ./bin -trimpath ./...

.PHONY: update
update:
	go get -u ./...
	go mod tidy

# Generate code including sqlc, wire
.PHONY: gen
gen: wire
	go generate ./...

.PHONY: sqlc
sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate -f ./internal/pkg/db/sqlc.yaml

# Run wire to generate wire_gen.go if needed
.PHONY: wire
wire:
	@echo "Running wire"
	go run github.com/google/wire/cmd/wire@latest ./...
