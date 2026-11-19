.PHONY: generate run clean test test-db

generate:
	sqlc generate
	templ generate

clean:
	docker compose down -v

run: clean generate
	@trap 'docker compose down -v; echo ""; echo "Servidor detenido correctamente."; exit 0' INT; \
	trap 'docker compose down -v' EXIT TERM; \
	docker compose up -d postgresql && \
	until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1 && \
	      docker exec postgres_container psql -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 1; \
	done && \
	go run main.go

test: clean generate
	@trap 'docker compose down -v' EXIT INT TERM; \
	docker compose up -d postgresql && \
	until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1 && \
	      docker exec postgres_container psql -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 1; \
	done && \
	go test -v ./...

test-db: clean generate
	@trap 'docker compose down -v' EXIT INT TERM; \
	docker compose up -d postgresql && \
	until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1 && \
	      docker exec postgres_container psql -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 1; \
	done && \
	go test -v ./db/sqlc
