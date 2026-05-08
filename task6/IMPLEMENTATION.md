# Реализация задач ДЗ №6

## 1. Интеграция Redis

Использована библиотека `github.com/redis/go-redis/v9`.
Создан платформенный пакет `platform/redis`, который инкапсулирует работу с Redis-клиентом.

Конфигурация находится в `api/config/config.yaml` и окружении:
```yaml
redis:
  host: "localhost"
  port: "6379"
```

## 2. Кэширование запросов

Реализовано через `CacheMiddleware` в `api/internal/controller/http/middleware.go`.

**Стратегия формирования ключей:**
Ключ формируется как `cache:<Path>:<RawQuery>`.
Например: `cache:/api/repositories/info:url=https://github.com/octocat/Hello-World`.

**Алгоритм:**
1. Для GET-запросов проверяется наличие ключа в Redis.
2. При `cache hit` — возвращается сохраненный JSON (данные не уходят к downstream сервисам).
3. При `cache miss` — выполняется бизнес-логика, ответ сохраняется в Redis.

**Порядок обработки:** Rate Limit -> Cache -> Handler.

## 3. TTL кэша

TTL задается в конфиге `cache.ttl_seconds` (по умолчанию 60 секунд). Применяется ко всем записям через `redis.Client.Set` с `cacheTTL`.

## 4. Rate Limiting

Реализован через `RateLimitMiddleware` в `api/internal/controller/http/middleware.go`.

**Описание реализации rate limiter:**
Используется алгоритм **Token Bucket** через библиотеку `golang.org/x/time/rate`.
Это **In-Memory** реализация (как разрешенный "глобальный limiter как упрощённый вариант" в ТЗ).

**Почему In-Memory:**
- ТЗ допускает упрощенный вариант.
- Используется `sync.Map` для хранения лимитеров по IP клиента.
- Если Redis недоступен, система корректно продолжает работать с In-Memory лимитером.

**Конфигурация:**
```yaml
rate_limit:
  requests_per_second: 5
  burst: 10
```

## 5. Обработка отказов

- Если Redis недоступен при старте (`Ping` не проходит), `redisClient` устанавливается в `nil`.
- Кэш при `nil` клиенте не используется (запросы идут напрямую к сервисам).
- Rate Limiter (In-Memory) продолжает работать независимо от Redis.

## 6. Docker

В `compose.yaml` добавлен сервис `redis:7-alpine`. 
API Gateway связывается с ним через переменную окружения `REDIS_HOST=redis`.

## Скриншоты (примеры логов)

1. **Cache Miss:**  
   `level=INFO msg="cache miss" key="cache:/api/repositories/info:url=..."`
2. **Cache Hit:**  
   `level=INFO msg="cache hit" key="cache:/api/repositories/info:url=..."`
3. **Rate Limit Exceeded:**  
   `level=WARN msg="rate limit exceeded" ip=...` и ответ API `429 Too Many Requests`.
