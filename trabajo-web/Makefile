.PHONY: generate run test test-db clean

generate:
	sqlc generate
	templ generate

clean:
	docker compose down -v

run: clean generate
	docker compose up -d postgresql
	@until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1 && \
	       docker exec postgres_container psql -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 1; \
	done
	go run main.go

test: clean generate
	docker compose up -d postgresql
	@until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1 && \
	       docker exec postgres_container psql -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 1; \
	done
	go test -v ./...
	docker compose down -v

test-db: clean generate
	docker compose up -d postgresql
	@until docker exec postgres_container pg_isready -U postgres -d outfits_db > /dev/null 2>&1 && \
	       docker exec postgres_container psql -U postgres -d outfits_db -c "SELECT 1;" > /dev/null 2>&1; do \
		sleep 1; \
	done
	go test -v ./db/sqlc
	docker compose down -v
