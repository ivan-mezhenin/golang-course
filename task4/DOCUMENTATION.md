## Запуск
```bash
docker compose build --no-cache
docker compose up -d
```
## Применение миграций
```bash
cd repo-stat/subscriber
migrate create -ext sql -dir migrations -seq create_subscriptions_table
```
## Примеры запросов

### Добавить новую подписку на репозиторий
```bash
curl -X POST http://localhost:28080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{"owner": "ivan-mezhenin", "repo": "golang-course"}' 
```

### Удалить подписку
```bash
curl -X DELETE http://localhost:28080/subscriptions/ivan-mezhenin/golang-course
```
### Получить список всех подписок
```bash
curl -X GET http://localhost:28080/subscriptions
```

### Получить информацию о всех репозиториях, на которых оформлена подписка
```bash
curl -X GET http://localhost:28080/subscriptions/info
```
