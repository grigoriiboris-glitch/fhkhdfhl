up:
	docker-compose up -d

b:
	docker-compose build
d:
	docker-compose down

app:
	docker-compose up api-backend
appb:
	docker-compose up --build api-backend
f:
	docker-compose up mindmap-frontend

m:
	docker-compose up migrate
c:
	docker-compose up caddy

test:
	docker exec -it mindmap-api go test ./...

td:
	docker exec -it -w /app/$(DIR) mindmap-api go test -cover

cr:
	docker exec -it  mindmap-api go run scripts/main.go
##	docker exec -it -w /app/auth mindmap-api go test -cover

# Значение по умолчанию
DIR ?= auth