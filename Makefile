# ==============================================================================
# 1. SYSTEM CONFIGURATION & ENVIRONMENT VARIABLES
# ==============================================================================

# Ensure Bash shell is used to avoid syntax errors
SHELL := /bin/bash

# --- Đọc các biến môi trường trực tiếp từ file .env ---
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# --- Database Connection String (Lấy thẳng từ biến môi trường) ---
DB_DSN := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

# Default env for Atlas (defaults to local if not provided)
ATLAS_ENV ?= local

# Find all .proto and .api files in the app/ directory
PROTO_FILES = $(shell find app -name "*.proto" -not -path "*/google/*")
API_FILES   = $(shell find app -name "*.api")

# ==============================================================================
# 2. UTILITY COMMANDS (HELPERS)
# ==============================================================================

.PHONY: help
help: ## Display list of commands
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: debug
debug: ## Check if environment variables are loaded correctly
	@echo "--- Debug Config ---"
	@echo "DB Host    : [$(DB_HOST)]"
	@echo "DB User    : [$(DB_USER)]"
	@echo "DB Name    : [$(DB_NAME)]"
	@echo "SSL Mode   : [$(DB_SSL_MODE)]"
	@echo "Full DSN   : $(DB_DSN)"
	@echo "Atlas Env  : [$(ATLAS_ENV)]"
	@echo "PayPal Mode: [$(PAYPAL_MODE)]"

# ==============================================================================
# 3. DATABASE MANAGEMENT WITH ATLAS (MIGRATIONS)
# ==============================================================================
# Note: SQL file directories are already configured in atlas.hcl
# Usage: make diff name=add_stock ATLAS_ENV=product

.PHONY: diff
diff: ## Create a new migration file from Go code. Ex: make diff name=add_users ATLAS_ENV=user
	@if [ -z "$(name)" ]; then echo "Error: Missing migration name. Please add name=migration_name"; exit 1; fi
	atlas migrate diff $(name) --env $(ATLAS_ENV)

.PHONY: apply
apply: ## Apply SQL files to the actual Database
	atlas migrate apply --env $(ATLAS_ENV) --url "$(DB_DSN)"

.PHONY: down
down: ## Rollback DB. Ex: make down OR make down v=20260104 ATLAS_ENV=product
	@if [ -n "$(v)" ]; then \
		echo "Reverting to version: $(v) for service $(ATLAS_ENV)..."; \
		atlas migrate down --env $(ATLAS_ENV) --url "$(DB_DSN)" --to-version "$(v)"; \
	else \
		echo "Reverting the latest migration step for service $(ATLAS_ENV)..."; \
		atlas migrate down --env $(ATLAS_ENV) --url "$(DB_DSN)"; \
	fi

.PHONY: migrate-hash
migrate-hash: ## Update hash if you manually edited .sql files
	atlas migrate hash --env $(ATLAS_ENV)

.PHONY: status
status: ## Check migration status in the DB
	atlas migrate status --env $(ATLAS_ENV) --url "$(DB_DSN)"

# ==============================================================================
# 4. CODE GENERATION & RUN SERVICES
# ==============================================================================

# Khai báo các lệnh giả (không phải tên file)
.PHONY: gen rpc gw

# Lệnh sinh mã nguồn Protobuf, Validate và Gateway Descriptor
gen:
	@echo "1. Đang sinh mã protobuf validation..."
	protoc -I . \
		--validate_out="lang=go,paths=source_relative:./dropshipbe" \
		dropshipbe.proto
		
	@echo "2. Đang sinh mã nguồn gRPC (go-zero)..."
	goctl rpc protoc dropshipbe.proto --go_out=. --go-grpc_out=. --zrpc_out=.
	
	@echo "3. Đang sinh tệp mô tả cho Gateway..."
	protoc -I . --include_imports --descriptor_set_out=dropshipbe.pb dropshipbe.proto
	
	@echo "Hoàn tất sinh mã!"

# Lệnh chạy máy chủ gRPC
rpc:
	@echo "Khởi động máy chủ gRPC..."
	go run dropshipbe.go -f etc/dropshipbe.yaml

# Lệnh chạy máy chủ Gateway
gw:
	@echo "Khởi động API Gateway..."
	go run gateway/gateway.go -f etc/gateway.yaml