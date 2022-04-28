build:
	docker-compose --env-file .env -f deployment/docker-compose.yml build

up:
	docker-compose --env-file .env -f deployment/docker-compose.yml up -d

down:
	docker-compose --env-file .env -f deployment/docker-compose.yml down