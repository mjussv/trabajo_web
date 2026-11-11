.PHONY: generate run test test-db

generate:
	sqlc generate
	templ generate

run: generate
	docker compose up -d postgresql
	@until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1 && \
	       docker exec postgres_container psql -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 1; \
	done
	go run main.go

# Ejecuta TODOS los tests del proyecto
test: generate
	docker compose down -v
	docker compose up -d postgresql
	@until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1 && \
	       docker exec postgres_container psql -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 1; \
	done
	go test -v ./...
	docker compose down -v

# Ejecuta SOLO los tests de la base de datos (db/sqlc)
test-db: generate
	docker compose down -v
	docker compose up -d postgresql
	@until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1 && \
	       docker exec postgres_container psql -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 1; \
	done
	go test -v ./db/sqlc
	docker compose down -v
