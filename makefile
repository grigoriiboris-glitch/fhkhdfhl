up:
	docker-compose up

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