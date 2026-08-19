# Ручное тестирование API

Для каждого из 5 обязательных эндпоинтов приведён один корректный и один ошибочный сценарий. Тесты выполнены последовательно на чистом (только что запущенном) сервере — `go run cmd/server/main.go`, порт `8080`.

---

## 1. `POST /tasks` — создание задачи

### ✅ Корректный сценарий

```bash
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Купить молоко"}'
```

![POST /tasks — успешное создание](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/01-post-tasks-success.png)

Результат: `201 Created` — задача создана с `id=1`, сгенерированным `created_at`, `done=false`.

### ❌ Ошибочный сценарий (пустой title)

```bash
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": ""}'
```

![POST /tasks — пустой title, 400](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/02-post-tasks-empty-title.png)

Результат: `400 Bad Request` — `{"error":"title is required"}`.

---

## 2. `GET /tasks` — список задач

### ✅ Корректный сценарий

```bash
curl -i http://localhost:8080/tasks
```

![GET /tasks — список задач](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/03-get-tasks-list.png)

Результат: `200 OK` — массив с одной ранее созданной задачей.

### ❌ Ошибочный сценарий (неподдерживаемый метод)

```bash
curl -i -X PATCH http://localhost:8080/tasks
```

![PATCH /tasks — 405 Method Not Allowed](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/04-get-tasks-method-not-allowed.png)

Результат: `405 Method Not Allowed` — `{"error":"method not allowed"}`.

---

## 3. `GET /tasks/{id}` — получение задачи по ID

### ✅ Корректный сценарий

```bash
curl -i http://localhost:8080/tasks/1
```

![GET /tasks/1 — успех](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/05-get-task-by-id.png)

Результат: `200 OK` — данные задачи с `id=1`.

### ❌ Ошибочный сценарий (несуществующий ID)

```bash
curl -i http://localhost:8080/tasks/999
```

![GET /tasks/999 — 404 Not Found](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/06-get-task-not-found.png)

Результат: `404 Not Found` — `{"error":"task not found"}`.

---

## 4. `PUT /tasks/{id}` — обновление задачи

### ✅ Корректный сценарий

```bash
curl -i -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Купить молоко и хлеб", "done": true}'
```

![PUT /tasks/1 — успешное обновление](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/07-put-task-update.png)

Результат: `200 OK` — `title` и `done` обновлены, `id` и `created_at` сохранены от исходной задачи.

### ❌ Ошибочный сценарий (несуществующий ID)

```bash
curl -i -X PUT http://localhost:8080/tasks/999 \
  -H "Content-Type: application/json" \
  -d '{"title": "Не существует"}'
```

![PUT /tasks/999 — 404 Not Found](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/08-put-task-not-found.png)

Результат: `404 Not Found` — `{"error":"task not found"}`.

---

## 5. `DELETE /tasks/{id}` — удаление задачи

### ✅ Корректный сценарий

```bash
curl -i -X DELETE http://localhost:8080/tasks/1
```

![DELETE /tasks/1 — 204 No Content](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/09-delete-task-success.png)

Результат: `204 No Content` — тело ответа пустое, это ожидаемо для этого кода.

### ❌ Ошибочный сценарий (повторное удаление той же задачи)

```bash
curl -i -X DELETE http://localhost:8080/tasks/1
```

![Повторный DELETE /tasks/1 — 404 Not Found](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/10-delete-task-not-found.png)

Результат: `404 Not Found` — задача уже удалена шагом выше.

---

## Дополнительно: `GET /health`

```bash
curl -i http://localhost:8080/health
```

![GET /health — 200 OK](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/11-health-check.png)

Результат: `200 OK` — `{"status":"ok"}`.

---

## Дополнительно: невалидный ID в пути

```bash
curl -i http://localhost:8080/tasks/abc
```

![GET /tasks/abc — 400 невалидный id](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/12-invalid-task-id.png)

Результат: `400 Bad Request` — `{"error":"invalid task id"}`.

---

## Дополнительно: некорректный JSON в теле запроса

```bash
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{invalid json'
```

![POST /tasks — невалидный JSON, 400](https://github.com/nikolaykonkin/go-tasks-api-crud-hw/blob/main/img/13-invalid-json-body.png)

Результат: `400 Bad Request` — `{"error":"invalid JSON body"}`.
