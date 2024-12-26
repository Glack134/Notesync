#Start BD --> docker run --name=notesync-db -e POSTGRES_PASSWORD='qwe230405180405' -p 5433:5432 -d --rm postgres

#Start app
.PHONY: run
run:
	go run cmd/main.go

.DEFAULT_GOAL := run

#Запуск докера
.PHONY: build up down

build:
	docker-compose build

up:
	docker-compose up

down:
	docker-compose down

