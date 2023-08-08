# Определите имя вашего исполняемого файла
APP_NAME = todo

# Цель для запуска контейнера MongoDB с параметрами
mongo:
	docker run --name region-todo_mongo_1 -p 2717:27017 -e MONGO_INITDB_ROOT_USERNAME=mongo -e MONGO_INITDB_ROOT_PASSWORD=123 -d mongo:4.4.23-focal

# Цель для запуска оболочки в Docker-контейнере
mng:
	docker exec -it region-todo_mongo_1 bash

# Цель для сборки Docker-образа
docker-build:
	docker build -t $(APP_NAME) .

# Цель для запуска Docker-compose
docker-up:
	docker-compose up

# Цель для остановки и удаления контейнеров
docker-down:
	docker-compose down

# Цель для остановки и удаление вместе с данными в базе данных контейнеров
docker-volume-down:
	docker stop region-todo_mongo_1
	docker rm region-todo_mongo_1
	docker run --rm -v mongodb-data:/data/db mongo:4.4.23-focal rm -rf /data/db/*
	docker-compose down
