ifneq (,$(wildcard ./.env))
    include .env
    export
endif


.PHONY: migrate_db
migrate_db:
	flyway -user=${POSTGRES_USER} -password=${POSTGRES_PASSWORD} -url=jdbc:${JDBC_URL} -locations=filesystem:sql -outputFile=/dev/null migrate


.PHONY: flyway
flyway:
	flyway -user=${POSTGRES_USER} -password=${POSTGRES_PASSWORD} -url=jdbc:${JDBC_URL} -locations=filesystem:sql -outputFile=/dev/null $(COMMAND)


.PHONY: flush_db
flush_db:
	echo "This command cleans up the whole database. Press ENTER if you are sure you are not on production, else press Ctrl+C"; \
	read REPLY; \
	docker-compose stop postgres; \
	rm -rf data/postgres || true; \
	docker-compose  up -d postgres; \
	until pg_isready -h localhost -p ${POSTGRES_PORT}; do sleep 1; done; \


.PHONY: fresh_db
fresh_db: flush_db migrate_db


.PHONY: rebuild
rebuild:
	docker builder prune -f
	docker-compose build --force-rm --no-cache


.PHONY: up
up:
ifeq ($(DEBUG), true)
	@docker-compose up -d redis postgres nats
else
	@migrate_db
	@rebuild
	@docker-compose up -d
endif

docs:
	cd src/services/api && swag init --parseDependency --parseInternal

