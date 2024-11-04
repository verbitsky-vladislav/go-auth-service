# Запуск проекта

## Шаги для запуска

### 1. Сгенерировать документацию
Чтобы сгенерировать документацию API с использованием Swagger, выполните следующую команду:
```shell
swag init -g cmd/main.go
```

### 2. Запустить Redis и PostgreSQL
Для запуска Redis и PostgreSQL используйте Docker Compose. Выполните следующую команду:
```shell
sudo docker-compose up --build
```

### 3. Запустить миграции в PostgreSQL
Убедитесь, что миграции инициализированы в PostgreSQL. Выполните необходимые команды миграции (например, с использованием миграционных инструментов, если они у вас есть).

### 4. Запустить проект
Для запуска вашего приложения выполните команду:
```shell
go run cmd/main.go
```

## Примечание
Убедитесь, что все зависимости установлены и что вы находитесь в корневом каталоге проекта перед выполнением этих команд.

---

# API Documentation

## Authentication API

### 1. Register User
- **Endpoint**: `/api/auth/register`
- **Method**: `POST`
- **Description**: Регистрация нового пользователя.
- **Request Body**:
    - `username` (string): Имя пользователя (обязательное).
    - `email` (string): Электронная почта (обязательная).
    - `password` (string): Пароль (обязательный).
- **Responses**:
    - **201 Created**: Успешная регистрация.
    - **400 Bad Request**: Ошибка валидации (например, отсутствует обязательное поле).

---

### 2. Login User
- **Endpoint**: `/api/auth/login`
- **Method**: `POST`
- **Description**: Аутентификация пользователя.
- **Request Body**:
    - `email` (string): Электронная почта (обязательная).
    - `password` (string): Пароль (обязательный).
- **Responses**:
    - **200 OK**: Успешный вход, возвращает токен аутентификации.
    - **401 Unauthorized**: Неверные учетные данные.

---

### 3. Logout User
- **Endpoint**: `/api/auth/logout`
- **Method**: `POST`
- **Description**: Выход пользователя, недействителен токен сессии.
- **Responses**:
    - **200 OK**: Успешный выход.
    - **401 Unauthorized**: Пользователь не аутентифицирован.

---

### 4. Get User Details
- **Endpoint**: `/api/auth/my`
- **Method**: `GET`
- **Description**: Получение информации о текущем пользователе.
- **Responses**:
    - **200 OK**: Возвращает данные пользователя.
    - **401 Unauthorized**: Пользователь не аутентифицирован.

---

### 5. Confirm Email Verification
- **Endpoint**: `/api/auth/confirm/verify/:verification_token`
- **Method**: `GET`
- **Description**: Подтверждение электронной почты с использованием токена подтверждения.
- **Path Parameters**:
    - `verification_token` (string): Токен подтверждения (обязательный).
- **Responses**:
    - **302 Found**: Успешное подтверждение, перенаправление на URL фронтенда.
    - **400 Bad Request**: Неверный токен подтверждения.

---

### 6. Google Login
- **Endpoint**: `/api/social/google/login`
- **Method**: `GET`
- **Description**: Перенаправляет пользователя на страницу входа через Google.
- **Responses**:
    - **302 Found**: Перенаправление на страницу аутентификации Google.

---

### 7. Google Callback
- **Endpoint**: `/api/social/google/callback`
- **Method**: `GET`
- **Description**: Обработка обратного вызова после аутентификации через Google.
- **Responses**:
    - **200 OK**: Успешная аутентификация и получение данных пользователя.

---

### 8. Yandex Login
- **Endpoint**: `/api/social/yandex/login`
- **Method**: `GET`
- **Description**: Перенаправляет пользователя на страницу входа через Yandex.
- **Responses**:
    - **302 Found**: Перенаправление на страницу аутентификации Yandex.

---

### 9. Yandex Callback
- **Endpoint**: `/api/social/yandex/callback`
- **Method**: `GET`
- **Description**: Обработка обратного вызова после аутентификации через Yandex.
- **Responses**:
    - **200 OK**: Успешная аутентификация и получение данных пользователя.

---

## Conclusion
Это документация по API для микросервиса аутентификации. Каждый эндпоинт описан с указанием метода, пути, параметров и возможных ответов. Используйте эти эндпоинты для интеграции с системой аутентификации.

