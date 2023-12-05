.DEFAULT_GOAL := up

up:
	docker-compose up jayplus-backend --build

up-detached:
	docker-compose up -d jayplus-backend --build

down:
	docker-compose down && docker-compose rm -f

