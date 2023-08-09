# Todo List Microservice

Этот микросервис предоставляет API для управления задачами Todo List.

## Требования

- Go (версия >= 1.15)
- Docker
- Docker Compose

## Установка

1. Клонируйте репозиторий:

```bash
git clone https://github.com/ctuzelov/region-todo.git
cd region-todo
```


## Команды с Makefile для запуска


### Запускать и останавливать микросервис с помощью Docker Compose
```bash
make docker-up
make docker-down
```
### Открыть оболочку Docker-контейнере
```bash
make mng
```

### запуск unit-тестов
```bash
make test
```

# Использование

### Открытие swagger документаций по ссылке после запуска

```bash
http://localhost:8080/swagger/index.html
```

### Создание новой задачи


    URL: POST /api/todo-list/tasks

Body:

```json
{
   "title": "Купить книгу",
   "activeAt": "2023-08-04"
}
```

### Обновление существующей задачи

    URL: PUT /api/todo-list/tasks/{ID}

Body:

```json
{
   "title": "Купить книгу - Высоконагруженные приложения",
   "activeAt": "2023-08-05"
}
```

### Удаление задачи

    URL: DELETE /api/todo-list/tasks/{ID}

### Пометить задачу выполненной

    URL: PUT /api/todo-list/tasks/{ID}/done

### Список задач по статусу

    URL: GET /api/todo-list/tasks

Query параметр: status (возможные значения: active или done), по умолчанию active


## Author
- [@ctuzelov](https://github.com/ctuzelov)

