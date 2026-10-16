# indica a la herramienta make que "test", "test-db" y "test-api" son tareas y no archivos
.PHONY: test test-db test-api

# ejecuta todos los tests que encuentre 
test:
	docker compose down -v
	docker compose up -d postgresql
	@until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1 && \
	       docker exec postgres_container psql -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 2; \
	done
	sqlc generate
	go test -v ./...
	docker compose down -v

# ejecuta SOLO la 2da entrega (tests de la BD en db/sqlc)
test-db:
	docker compose down -v
	docker compose up -d postgresql
	@until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1 && \
	       docker exec postgres_container psql -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 2; \
	done
	sqlc generate
	go test -v ./db/sqlc
	docker compose down -v

# ejecuta SOLO la 3ra entrega (tests de la API en la raíz)
test-api:
	docker compose down -v
	docker compose up -d postgresql
	@until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1 && \
	       docker exec postgres_container psql -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 2; \
	done
	sqlc generate
	go test -v .
	docker compose down -v
