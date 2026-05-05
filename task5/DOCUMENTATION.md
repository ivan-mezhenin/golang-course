## Топики Kafka
`repo-request` 
* Назначение: Processor отправляет запросы на сбор информации о репозиториях
```json
{
  "owner": "golang",
  "repo": "go"
}
```
`repo-responses` 
* Назначение: Collector отправляет результаты сбора данных
```json
{
  "owner": "golang",
  "repo": "go",
  "full_name": "golang/go",
  "description": "The Go programming language",
  "stars": 125000,
  "forks": 18000,
  "created_at": "2009-11-10T23:27:04Z",
  "visibility": "public",
  "error": "" 
}
```
`subscription-updates` 
* Назначение: Collector каждые 15 секунд запрашивает список подписок у Subscribe и публикует их в этот топик. Processor читает и обновляет локальную таблицу подписок.
```json
[
  {"owner": "golang", "repo": "go"},
  {"owner": "kubernetes", "repo": "kubernetes"}
]
```

## Запуск
```bash
docker compose build --no-cache
docker compose up -d
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
