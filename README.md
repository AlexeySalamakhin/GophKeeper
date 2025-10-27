# GophKeeper

### Сборка проекта

go mod download
go mod tidy

### Тесты

go test ./...

# Сборка сервера и клиента

go build -o build/gophkeeper-server ./cmd/server
go build -o build/gophkeeper-client ./cmd/client

## Запуск

### Сервер

./build/gophkeeper-server

```

Сервер будет доступен по адресу `http://localhost:8080`

### Клиент

# Или напрямую
./build/gophkeeper-client
```

## Использование

### Регистрация пользователя

./build/gophkeeper-client auth register username email@example.com password

### Вход в систему

./build/gophkeeper-client auth login username password

### Работа с данными

# Список всех данных

./build/gophkeeper-client data list

# Список данных определенного типа

./build/gophkeeper-client data list login_password

# Добавление новых данных

./build/gophkeeper-client data add login_password "Мой сайт"

# Получение данных по ID

./build/gophkeeper-client data get <id>

### Проверка версии

./build/gophkeeper-client version

## API Endpoints

### Аутентификация

- `POST /api/v1/register` - Регистрация пользователя
- `POST /api/v1/login` - Вход в систему

### Данные (требуют авторизации)

- `GET /api/v1/data` - Получение всех данных пользователя
- `GET /api/v1/data/{id}` - Получение данных по ID
- `POST /api/v1/data` - Создание новых данных
- `PUT /api/v1/data/{id}` - Обновление данных
- `DELETE /api/v1/data/{id}` - Удаление данных
- `GET /health` - Проверка состояния сервера

### Переменные окружения

#### 1. Файл .env (рекомендуется для разработки)

Создайте файл `.env` в корне проекта:

SERVER_HOST=localhost
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=gophkeeper
DB_PASSWORD=password
DB_NAME=gophkeeper
DB_SSLMODE=disable
JWT_SECRET=your-secret-key-change-in-production

#### 2. Системные переменные окружения

- `SERVER_HOST` - хост сервера
- `SERVER_PORT` - порт сервера
- `DB_HOST` - хост PostgreSQL
- `DB_PORT` - порт PostgreSQL
- `DB_USER` - пользователь PostgreSQL
- `DB_PASSWORD` - пароль PostgreSQL
- `DB_NAME` - имя базы данных
- `DB_SSLMODE` - режим SSL
- `JWT_SECRET` - секретный ключ для JWT
