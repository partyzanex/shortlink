POSTGRES_DSN ?= postgresql://postgres:postgres@localhost:5432/short?sslmode=disable
GENERATED_GO_PATH = $(CURDIR)/pkg/proto/go

include proto.mk

proto.mk:
	@tmpdir=$$(mktemp -d) && \
	git clone --depth 1 --single-branch https://github.com/partyzanex/go-makefile.git $$tmpdir && \
	cp $$tmpdir/proto.mk $(CURDIR)/proto.mk

CLI_CONFIG_GEN_VERSION := v0.0.5
CLI_CONFIG_GEN_BIN := $(LOCAL_BIN)/cli-config-gen

.PHONY: config
config: cli-config-gen-install
	@$(CLI_CONFIG_GEN_BIN) -s ./app/config.yaml -t ./internal/config/config.go \
	&& go fmt ./internal/config

MIGRATIONS_PATH := $(CURDIR)/migrations

.PHONY: migration
migration: goose-install ## Creating migrations
ifneq ($(wildcard $(MIGRATIONS_PATH)),)
	@read -p "Enter migration name: " migration_name; \
	$(GOOSE_BIN) -dir $(MIGRATIONS_PATH) create $$migration_name sql
endif

.PHONY: up-env
up-env: up-down
	# @docker-compose pull
	@docker-compose up -d

.PHONY: up-down
up-down:
	@docker-compose down

SQLBOILER_VERSION=v4.14.2
SQLBOILER_BIN=$(LOCAL_BIN)/sqlboiler
SQLBOILER_DRIVER_BIN=$(LOCAL_BIN)/sqlboiler-psql

.PHONY: sqlboiler-install
sqlboiler-install:
	@go-install github.com/volatiletech/sqlboiler/v4@$(SQLBOILER_VERSION) $(SQLBOILER_BIN)
	@go-install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@$(SQLBOILER_VERSION) $(SQLBOILER_DRIVER_BIN)

MAKE_PATH=$(LOCAL_BIN):/bin:/usr/bin:/usr/local/bin

.PHONY: sqlboiler-gen
sqlboiler-gen: sqlboiler-install up-env pg-wait goose-up
	PATH=$(MAKE_PATH) $(SQLBOILER_BIN) psql

.PHONY: build
build:
	go build -o ./bin/shortlink ./cmd/shortlink