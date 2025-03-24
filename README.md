# Tarantool KV Storage API

Этот репозиторий содержит реализацию HTTP-сервера на Go, который предоставляет RESTful API для работы с Tarantool в качестве key-value хранилища.

## Инструкция по запуску

1. Клонируйте репозиторий:
   
   ```bash
   git clone https://github.com/p0rtale/tarantool-kv.git
   cd tarantool-kv
   ```
   
2. Настройте переменные окружения:

   ```
   TARANTOOL_HOST=<имя_хоста_tarantool>
   TARANTOOL_PORT=<порт_tarantool>
   TARANTOOL_USER=<имя_пользователя_tarantool>
   TARANTOOL_PASSWORD=<пароль_tarantool>

   GRAFANA_USER=<имя_пользователя_grafana>
   GRAFANA_PASSWORD=<пароль_grafana>
   ```
   
3. Запустите контейнеры:

   ```
   docker-compose up --build -d
   ```

## API-документация 

### POST /kv

Создает новую запись.

* Тело запроса:

   ```json
   {
     "key": "test",
     "value": { "name": "Alice" }
   }
   ```

* Возможные ответы:
  
  * `201 Created` — запись успешно создана.
  * `409 Conflict` — ключ уже существует.
  * `400 Bad Request` — некорректное тело запроса.

### PUT /kv/{id}

* Тело запроса:

   ```json
   {
     "value": { "name": "Bob" }
   }
   ```

* Возможные ответы:
  
  * `200 OK` — запись успешно обновлена.
  * `404 Not Found` — ключ не существует.
  * `400 Bad Request` — некорректное тело запроса.

### GET /kv/{id} 

Получает значение по ключу. 

* Возможные ответы:
  
  * `200 OK` — запись успешно получена.
  * `404 Not Found` — ключ не существует.

### DELETE /kv/{id} 

Удаляет запись по ключу. 

* Возможные ответы:
  
  * `200 OK` — запись успешно удалена.
  * `404 Not Found` — ключ не существует.

## Логирование

Все операции логируются в консоль. Чтобы посмотреть логи, выполните команду: 

  ```
  docker logs kv-server
  ```

Формат логов:

  ```
  Method=POST URL=/kv RemoteAddr=127.0.0.1:54321 Status=201 Duration=0.123
  ```

## Мониторинг и метрики

Для мониторинга работы приложения используются Prometheus и Grafana (доступна на порту `3000`). 

### Доступные метрики:

1. `http_requests_total`
   
   * Общее количество HTTP-запросов.
   * Метки: `method`, `endpoint`, `status`. 

3. `http_request_duration_seconds`

   * Время обработки HTTP-запросов.
   * Метки: `method`, `endpoint`. 

