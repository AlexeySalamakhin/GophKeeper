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

### База данных PostgreSQL с Docker Compose

Проект использует Docker Compose для автоматического создания и управления базой данных PostgreSQL. База данных создается автоматически при первом запуске контейнера.

#### Быстрый старт

1. Запустите PostgreSQL контейнер:

```bash
docker-compose up -d
```

Это автоматически:

- Создаст базу данных `gophkeeper`
- Создаст пользователя `gophkeeper` с паролем `password`
- Настроит подключение на порту `5433`
- Сохранит данные в Docker volume `postgres_data`

#### Управление контейнером

**Остановка базы данных** (данные сохраняются):

```bash
docker-compose stop
```

**Запуск остановленной базы данных**:

```bash
docker-compose start
```

**Остановка и удаление контейнера** (данные сохраняются в volume):

```bash
docker-compose down
```

**Остановка и удаление контейнера со всеми данными**:

```bash
docker-compose down -v
```

### Сервер

**Важно:** Перед запуском сервера убедитесь, что база данных PostgreSQL запущена через `docker-compose up -d`.

1. Убедитесь, что файл `.env` настроен правильно (см. раздел "Переменные окружения")

2. Запустите сервер:

```bash
./build/gophkeeper-server
```

Сервер будет доступен по адресу `http://localhost:8080`

### Клиент

# Или напрямую

./build/gophkeeper-client

````

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

Создайте файл `.env` в корне проекта на основе `env.example`:

```bash
cp env.example .env
````

Затем отредактируйте `.env` и установите необходимые значения. **Важно:** следующие переменные обязательны и не имеют значений по умолчанию:

- `JWT_SECRET` - секретный ключ для JWT (обязательно)
- `CRYPTO_KEY` - ключ шифрования (обязательно)
- `DB_PASSWORD` - пароль PostgreSQL (обязательно)
- `DB_USER` - пользователь PostgreSQL (обязательно)
- `DB_NAME` - имя базы данных (обязательно)

Пример `.env` для работы с docker-compose:

```env
SERVER_HOST=localhost
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=gophkeeper
DB_PASSWORD=password
DB_NAME=gophkeeper
DB_SSLMODE=disable
JWT_SECRET=your-secret-key-change-in-production
CRYPTO_KEY=your-encryption-key-change-in-production
```

#### 2. Системные переменные окружения

Все переменные окружения можно установить через системные переменные:

- `SERVER_HOST` - хост сервера (по умолчанию: localhost)
- `SERVER_PORT` - порт сервера (по умолчанию: 8080)
- `DB_HOST` - хост PostgreSQL (по умолчанию: localhost)
- `DB_PORT` - порт PostgreSQL (по умолчанию: 5432)
- `DB_USER` - пользователь PostgreSQL (**обязательно**)
- `DB_PASSWORD` - пароль PostgreSQL (**обязательно**)
- `DB_NAME` - имя базы данных (**обязательно**)
- `DB_SSLMODE` - режим SSL (по умолчанию: disable)
- `JWT_SECRET` - секретный ключ для JWT (**обязательно**)
- `CRYPTO_KEY` - ключ шифрования (**обязательно**)
