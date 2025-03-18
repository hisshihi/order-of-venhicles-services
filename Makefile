createdb:
	docker exec -it order-services createdb --username=root --owner=root order_services

dropdb:
	docker exec -it order-services dropdb --username=root order_services

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/order_services?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/order_services?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/order_services?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/order_services?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run cmd/main.go -config=cmd/app.env

.PHONY: createdb dropdb migrateup migratedown migratedown1 migrateup1 sqlc test server 