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
