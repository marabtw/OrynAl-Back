build:
	docker compose up --build -d

down:
	docker compose down

migrate_up:
	migrate -path ./app/schema -database 'postgres://postgres:postgres@localhost:5432/orynal_db?sslmode=disable' up

migrate_down:
	 migrate -path ./app/schema -database 'postgres://postgres:postgres@localhost:5432/orynal_db?sslmode=disable' down

migrate_admin:
	./insert_user.sh
