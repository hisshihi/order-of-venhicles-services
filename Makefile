createdb:
	docker exec -it order-services createdb --username=root --owner=root order_services

dropdb:
	docker exec -it order-services dropdb --username=root order_services

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/order_services?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/order_services?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: createdb dropdb migrateup migratedown