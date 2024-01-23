DB_URL=postgresql://root:secret@localhost:5432/simplebank?sslmode=disable

postgres:
	docker run --name simplebank-pg -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it simplebank-pg createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it simplebank-pg dropdb simplebank

migrateup:
	migrate --path db/migration --database "$(DB_URL)" --verbose up

migrateup1:
	migrate --path db/migration --database "$(DB_URL)" --verbose up 1

migratedown:
	migrate --path db/migration --database "$(DB_URL)" --verbose down

migratedown1:
	migrate --path db/migration --database "$(DB_URL)" --verbose down 1

dbdocs:
	dbdocs build doc/db.dbml

dbschema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/itsmetambui/simplebank/db/sqlc Store

server:
	go run main.go

server-docker:
	docker run --name simplebank --network bank-network -p 8080:8080 -e DB_SOURCE="postgresql://root:secret@simplebank-pg:5432/simplebank?sslmode=disable" simplebank:latest

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 dbdocs dbschema sqlc test server server-docker