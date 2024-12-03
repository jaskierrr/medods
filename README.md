# medods
# Стек
- стандартная http либа
- БД Postgres
- Миграции golang-migrate
- Сборка Docker
- Старался придерживаться Clean Architecture

# Сборка
`make build` в корневой папке

# EndPoints
есть две ручки:
- `curl -X POST http://127.0.0.1:8080/login?id=12311`
- `curl -X POST http://127.0.0.1:8080/refresh \
-H "Content-Type: application/json" \
-d '{
  "refresh_token": "JDJhJDEwJGdWLlFJNlNJakpxSmJMQTcwMER2dmUyT0V1WEtVcm9hUkl6eDI3UUxuM0VGQlJHQmxNODFT",
  "user": {
    "id": "12311"
  }
}
`

# Примечание
- все требования учтены
- при изменении ip адресса (достаточно перезапустить контейнер) в лог пишется сообщение о рефреше токенов (через мок email репозиторий)
- успел написать один тест, он довольно страшненький, я торопился
- время окончания работы токенов можно поменять в .env файле
