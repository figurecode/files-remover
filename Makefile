.PHONY: all

SHELL=/bin/bash -e

help: ## Справка
	@awk 'BEGIN {FS = ":.*?## "} /^[0-9-a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

info: ## Шпаргалка по установки из README.md
	@sed '/git/,/```/!d;/```/q' README.md | grep -v '```'

ps:
	docker ps | grep --color alg-go

rebuild: ## Сборка контейнеров без запуска проекта
	docker compose -f docker/docker-compose.yml build --no-cache

up: ## Старт контейнера
	docker compose -f docker/docker-compose.yml up -d 

down: ## Остановить контейнер
	docker compose -f docker/docker-compose.yml down

bash: ## Зайти в контейнер go
	docker compose -f docker/docker-compose.yml exec go sh
