.PHONY: build run updb runserver

build:
	 docker build -t db-forum .

run:
	docker run -it --rm -p 5432:5432 --name db-forum -t db-forum

runserver:
	go run ./cmd/db-forum/main.go


updb:
	docker-compose up

listen:
	lsof -i :5000 -P