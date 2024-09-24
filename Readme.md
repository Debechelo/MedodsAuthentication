# MedodsAuthentication

Проект запускается в Docker
```bash
docker-compose up --build
```

Приложение будет доступно по адресу http://localhost:8080


### Генерация токенов

- **URL**: `/token`
- **Метод**: `POST`
- **Параметры**:
    - `user_id`: ID пользователя, для которого генерируются токены.

**Пример запроса**:

```bash
curl -X POST "localhost:8080/token?user_id=test-user"
```

### Обновление токена
- **URL**: `/token/refresh`
- **Метод**: `POST`

**Пример запроса**:
```bash
curl -X POST "localhost:8080/token/refresh"
```

Тело запроса:
```json
{
"refresh_token": "refresh_token",
"user_id": "user_id"
}
```



# Тестирование
Для запуска тестов выполните команду:
```bash
go test ./tests -v
```