# Online Subscriptions Aggregation Service

REST-сервис для управления и агрегации данных об онлайн-подписках пользователей.

---

## Описание

Сервис позволяет:

- Создавать, читать, обновлять, удалять и листить записи о подписках (CRUDL).
- Считать суммарную стоимость подписок за выбранный период с фильтрацией по `user_id` и `service_name`.
- Хранить данные в PostgreSQL с поддержкой миграций.
- Использовать конфигурацию из `./config/configs.yaml` файла.
- Запускаться через Docker Compose.

---

##  Настройка

1. Заполните параметры:

```yaml
srv:
  port: "8080"
  
postgres:
  username: "postgres"
  password: "postgres"
  host: "localhost"
  port: 5432
  db_name: "postgres"```
```

2. Соберите и запустите сервис через Docker Compose:

```bash
docker-compose up --build
```

* Сервис будет доступен на `http://localhost:8080`.

---

## Миграции

* Используется PostgreSQL.
* Миграции создаются через GORM `AutoMigrate`.
* Таблица `subscriptions` создается с необходимыми полями и проверками (`price >= 0`, `start_date < end_date`).

---

## API Endpoints

### CRUDL для подписок

| Метод  | Путь                 | Описание                      |
| ------ | -------------------- | ----------------------------- |
| POST   | `/subscriptions`     | Создать новую подписку        |
| GET    | `/subscriptions`     | Получить список всех подписок |
| GET    | `/subscriptions/:id` | Получить подписку по ID       |
| PUT    | `/subscriptions/:id` | Обновить подписку             |
| DELETE | `/subscriptions/:id` | Удалить подписку              |

### Агрегация стоимости

| Метод | Путь                   | Параметры                                           |
| ----- | ---------------------- | --------------------------------------------------- |
| GET   | `/subscriptions/total` | `user_id`, `service_name`, `start_date`, `end_date` |

---

## Swagger документация

* Swagger документация находиться [ЗДЕСЬ](http://github.com/ashurov-imomali/sbscribtion-service/blob/main/Readme.md)

---

## Технологии

* Go (GORM, uuid)
* PostgreSQL
* Docker & Docker Compose
* Swagger =
