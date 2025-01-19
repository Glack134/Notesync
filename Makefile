#Start BD --> docker run --name=notesync-db -e POSTGRES_PASSWORD='qwe230405180405' -p 5432:5432 -d --rm postgres

#Start app
build:
	docker-compose build notesync

run:
	docker-compose up notesync

test:
	go test -v ./...

migrate:
	migrate -path ./schema -database 'postgres://midiy:qwe230405180405@0.0.0.0:5436/postgres?sslmode=disable' up

swag:
	swag init -g cmd/main.go

