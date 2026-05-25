up:
	docker compose -f deployments/compose/docker-compose.yml up -d

down:
	docker compose -f deployments/compose/docker-compose.yml down

logs:
	docker compose -f deployments/compose/docker-compose.yml logs -f

config:
	docker compose deployments/compose/docker-compose.yml config

build:
	docker compose -f deployments/compose/docker-compose.yml build

list:
	docker compose -f deployments/compose/docker-compose.yml ps