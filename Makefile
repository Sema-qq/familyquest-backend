MIGRATIONS_DIR := db/migrations

.PHONY: migrate-new openapi
migrate-new:
	@test -n "$(name)" || (echo "usage: make migrate-new name=create_user_table" && exit 1)
	@mkdir -p $(MIGRATIONS_DIR)
	@touch "$(MIGRATIONS_DIR)/$$(date +%Y%m%d%H%M%S)_$(name).sql"

openapi:
	@docker compose up swagger-ui
