mongo:
	docker run --name mongodb -p 2717:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=123 -d mongo