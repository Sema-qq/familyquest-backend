MIGRATIONS_DIR := db/migrations

.PHONY: migrate-new
migrate-new:
	@test -n "$(name)" || (echo "usage: make migrate-new name=create_user_table" && exit 1)
	@mkdir -p $(MIGRATIONS_DIR)
	@touch "$(MIGRATIONS_DIR)/$$(date +%Y%m%d%H%M%S)_$(name).sql"
