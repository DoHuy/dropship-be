# ==============================================================================
# 1. SYSTEM CONFIGURATION & ENVIRONMENT VARIABLES
# ==============================================================================

# Ensure Bash shell is used to avoid syntax errors
SHELL := /bin/bash

# --- Xác định file cấu hình YAML ---
# Bạn có thể thay đổi đường dẫn này trỏ tới file yaml thực tế của bạn
CONFIG_FILE := etc/dropshipbe.yaml

# Kiểm tra xem công cụ yq đã được cài đặt chưa
YQ_CHECK := $(shell command -v yq 2> /dev/null)
ifndef YQ_CHECK
$(error "Lỗi: Không tìm thấy 'yq'. Vui lòng cài đặt yq (https://github.com/mikefarah/yq) để đọc file YAML")
endif

# --- Parse cấu hình từ file YAML ---
# Sử dụng yq để lấy các giá trị nằm trong node Posgres
DB_HOST := $(shell yq '.Posgres.Host' $(CONFIG_FILE))
DB_PORT := $(shell yq '.Posgres.Port' $(CONFIG_FILE))
DB_USER := $(shell yq '.Posgres.User' $(CONFIG_FILE))
DB_PASS := $(shell yq '.Posgres.Password' $(CONFIG_FILE))
DB_NAME := $(shell yq '.Posgres.DBName' $(CONFIG_FILE))
DB_SSL  := $(shell yq '.Posgres.SSLMode' $(CONFIG_FILE))

# --- Database Connection String ---
DB_DSN := postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL)

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
	@echo "DB Host  : [$(DB_HOST)]"
	@echo "DB User  : [$(DB_USER)]"
	@echo "DB Name  : [$(DB_NAME)]"
	@echo "SSL Mode : [$(DB_SSL)]"
	@echo "Full DSN : $(DB_DSN)"
	@echo "Atlas Env: [$(ATLAS_ENV)]"

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



# Khai báo các lệnh giả (không phải tên file)
.PHONY: gen rpc gw

# Lệnh sinh mã nguồn Protobuf và Gateway Descriptor
gen:
	@echo "Đang sinh mã nguồn gRPC..."
	goctl rpc protoc dropshipbe.proto --go_out=. --go-grpc_out=. --zrpc_out=.
	@echo "Đang sinh tệp mô tả cho Gateway..."
	protoc --descriptor_set_out=dropshipbe.pb dropshipbe.proto
	@echo "Hoàn tất!"

# Lệnh chạy máy chủ gRPC
rpc:
	@echo "Khởi động máy chủ gRPC..."
	go run dropshipbe.go -f etc/dropshipbe.yaml

# Lệnh chạy máy chủ Gateway
gw:
	@echo "Khởi động API Gateway..."
	go run gateway/gateway.go -f etc/gateway.yaml	
