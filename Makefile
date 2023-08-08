mongo:
	docker run --name mongodb -p 2717:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=123 -d mongo:4.4.23-focal

mng:
	docker exec -it mongodb bash

# Определите имя вашего исполняемого файла
APP_NAME = todo

# Цель для сборки Docker-образа
docker-build:
	docker build -t $(APP_NAME) .

# Цель для запуска Docker-compose
docker-up:
	docker-compose up

# Цель для остановки и удаления контейнеров
docker-down:
	docker-compose down