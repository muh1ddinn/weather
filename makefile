
migrate-up:
	migrate -database "postgres://muhiddin:1@localhost:5432/weather?sslmode=disable" -path migration up

migrate-down:
	migrate -database "postgres://muhiddin:1@localhost:5432/weather?sslmode=disable" -path migration down

migrate-force:
	migrate -database "postgres://muhiddin:1@localhost:5432/weather?sslmode=disable" -path migration force 1
