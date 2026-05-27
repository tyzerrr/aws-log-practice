ifneq (,$(wildcard server/.env))
include server/.env
export
endif

BIN_DIR := $(abspath ./bin)
DB_PORT ?= 5432
DATABASE_URL ?= postgres://aws-log-practice:aws-log-practice@localhost:$(DB_PORT)/aws-log-practice?sslmode=disable
ATLAS := $(BIN_DIR)/atlas

.PHONY: setup
setup:
	@curl -sSf https://atlasgo.sh | sh -s -- -y --no-install -o $(CURDIR)/bin/atlas
	@chmod +x $(ATLAS)

.PHONY: db-up
db-up:
	@docker compose up -d --wait

.PHONY: db-down
db-down:
	@docker compose down -v

.PHONY: db-apply
db-apply: $(ATLAS)
	@[ "$(DATABASE_URL)" ] || ( echo ">> DATABASE_URL required"; exit 1 )
	@DATABASE_URL="$(DATABASE_URL)" $(ATLAS) migrate apply --env local

.PHONY: db-diff
db-diff: $(ATLAS)
	@[ "$(name)" ] || ( echo ">> name=<migration-name> required"; exit 1 )
	@$(ATLAS) migrate diff $(name) --env local

.PHONY: db-hash
db-hash: $(ATLAS)
	@$(ATLAS) migrate hash --env local

.PHONY: sqlc
sqlc:
	@go tool sqlc generate

.PHONY: buf-gen
buf-gen:
	@buf generate
