up:
	docker-compose up -d

b:
	docker-compose build
d:
	docker-compose down

app:
	docker-compose up api-backend
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
##	docker exec -it -w /app/auth mindmap-api go test -cover

# Значение по умолчанию
DIR ?= auth

#mkdir -p cmd/server internal/{auth,config,db,http/{requests/{auth,mindmap,post,user},middleware,routes,handlers},services,server} models repository/{migrations} storage/{uploads,cache,logs} pkg/{logger,utils} scripts