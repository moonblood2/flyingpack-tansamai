build: 
	docker-compose build
up:
	docker-compose up -d

down:
	docker-compose down

re:
	docker-compose down
	docker-compose up -d

run:
	docker-compose down
	docker-compose up -d
	docker logs -f an_server

sh:
	docker exec -it an_server sh

log:
	docker logs -f an_server