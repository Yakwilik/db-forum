.PHONY: build run updb runserver

build:
	 docker build -t db-forum .

run:
	docker run -it --memory 2G --log-opt max-size=5M --log-opt max-file=3 --rm -p 5432:5432  -p 5050:5050 -p 9001:9001 --name db-forum -t db-forum

runserver:
	go run ./cmd/db-forum/main.go


updb:
	docker-compose up

listen:
	lsof -i :5000 -P