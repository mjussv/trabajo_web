SHELL := bash
.PHONY: generate run run-app clean reset-db test test-db test-hurl

SESSION_SECRET ?= dev-secret-ootd-2026

generate:
	sqlc generate
	templ generate

clean:
	docker compose down -v

reset-db: clean

run: clean generate
	@trap 'docker compose down -v; echo ""; echo "Servidor detenido correctamente."; exit 0' INT; \
	trap 'docker compose down -v' EXIT TERM; \
	docker compose up -d postgresql && \
	until docker exec postgres_container pg_isready -h 127.0.0.1 -p 5432 -U postgres -d outfits_db > /dev/null 2>&1 && \
	      docker exec postgres_container psql -h 127.0.0.1 -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 1; \
	done && \
	SESSION_SECRET=$(SESSION_SECRET) go run main.go

run-app: run

test: clean generate
	@trap 'docker compose down -v' EXIT INT TERM; \
	docker compose up -d postgresql && \
	until docker exec postgres_container pg_isready -h 127.0.0.1 -p 5432 -U postgres -d outfits_db > /dev/null 2>&1 && \
	      docker exec postgres_container psql -h 127.0.0.1 -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 1; \
	done && \
	go test -v ./...

test-db: clean generate
	@trap 'docker compose down -v' EXIT INT TERM; \
	docker compose up -d postgresql && \
	until docker exec postgres_container pg_isready -h 127.0.0.1 -p 5432 -U postgres -d outfits_db > /dev/null 2>&1 && \
	      docker exec postgres_container psql -h 127.0.0.1 -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 1; \
	done && \
	go test -v ./db/sqlc

test-hurl: clean generate
	@trap 'docker compose down -v' EXIT INT TERM; \
	docker compose up -d --build && \
	until docker exec app_container wget -qO- http://localhost:8080/ > /dev/null 2>&1; do \
		sleep 1; \
	done && \
	hurl --test requests.hurl
