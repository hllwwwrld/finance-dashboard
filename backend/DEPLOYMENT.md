# Инструкция по развёртыванию на чистой системе

## Предварительные требования

### 1. Установка Docker

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Или используйте официальную документацию:
# https://docs.docker.com/engine/install/
```

### 2. Установка Docker Compose

Docker Compose обычно устанавливается вместе с Docker Desktop, но если нужно установить отдельно:

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install docker-compose-plugin

# Или используйте официальную документацию:
# https://docs.docker.com/compose/install/
```

Проверка установки:
```bash
docker --version
docker compose version
```

## Развёртывание приложения

### Шаг 1: Клонирование проекта

```bash
git clone <your-repo-url>
cd finance-dashboard/backend
```

### Шаг 2: Создание файла с переменными окружения

```bash
cp env.example .env
```

Отредактируйте `.env` файл и укажите реальные значения:

```bash
# PostgreSQL Configuration
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=your_database
POSTGRES_PORT=5432

# Application Database Connection
DB_HOST=postgres  # или внешний хост БД, если используете внешнюю БД
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_secure_password
DB_NAME=your_database

# Важно для production: установите require или verify-full для запуска миграций
DB_SSLMODE=require  # или verify-full для production

# Backend Port
BACKEND_PORT=3000
```

**Важно для production:**
- Установите `DB_SSLMODE=require` (или `verify-full`) для запуска миграций автоматически
- Используйте надёжные пароли
- Если используете внешнюю БД, укажите правильный `DB_HOST`

### Шаг 3: Сборка и запуск

```bash
# Сборка образов
docker compose build

# Запуск в фоновом режиме
docker compose up -d

# Просмотр логов
docker compose logs -f

# Просмотр статуса
docker compose ps
```

### Шаг 4: Проверка работы

```bash
# Проверить, что контейнеры запущены
docker compose ps

# Проверить логи миграций
docker compose logs migrations

# Проверить логи backend
docker compose logs backend

# Проверить работу API (если порт открыт)
curl http://localhost:3000/ping
```

## Управление приложением

### Остановка

```bash
docker compose down
```

### Остановка с удалением volumes (⚠️ удалит данные БД)

```bash
docker compose down -v
```

### Перезапуск

```bash
docker compose restart
```

### Пересборка после изменений

```bash
docker compose build --no-cache
docker compose up -d
```

### Просмотр логов

```bash
# Все сервисы
docker compose logs -f

# Конкретный сервис
docker compose logs -f backend
docker compose logs -f migrations
docker compose logs -f postgres
```

## Настройка для production

### 1. Использование внешней базы данных

Если у вас есть внешняя БД (например, managed PostgreSQL), обновите `.env`:

```bash
DB_HOST=your-db-host.com
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSLMODE=require  # важно для внешней БД
```

И удалите сервис `postgres` из `docker-compose.yml`, или создайте `docker-compose.prod.yml`:

```yaml
# docker-compose.prod.yml
services:
  migrations:
    # ... существующая конфигурация
  backend:
    # ... существующая конфигурация
    depends_on:
      migrations:
        condition: service_completed_successfully
    # убрать зависимость от postgres, если используете внешнюю БД
```

Запуск:
```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### 2. Настройка reverse proxy (nginx)

Пример конфигурации nginx:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. Автозапуск при перезагрузке сервера

Создайте systemd сервис `/etc/systemd/system/finance-dashboard.service`:

```ini
[Unit]
Description=Finance Dashboard Backend
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/path/to/finance-dashboard/backend
ExecStart=/usr/bin/docker compose up -d
ExecStop=/usr/bin/docker compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
```

Активация:
```bash
sudo systemctl enable finance-dashboard
sudo systemctl start finance-dashboard
```

## Устранение проблем

### Миграции не выполняются

Проверьте `DB_SSLMODE` в `.env`:
- Для локальной разработки: `DB_SSLMODE=disable` (миграции пропускаются)
- Для production: `DB_SSLMODE=require` или `verify-full` (миграции выполняются)

### Ошибки подключения к БД

1. Проверьте, что БД доступна:
   ```bash
   docker compose logs postgres
   ```

2. Проверьте credentials в `.env`

3. Если используете внешнюю БД, проверьте firewall и сетевые настройки

### Порт уже занят

Измените `BACKEND_PORT` в `.env` или освободите порт:

```bash
# Найти процесс на порту 3000
sudo lsof -i :3000

# Убить процесс
sudo kill -9 <PID>
```

## Безопасность

1. **Никогда не коммитьте `.env` файл в git**
2. Используйте сильные пароли в production
3. Настройте firewall для ограничения доступа к портам
4. Используйте SSL/TLS для подключения к БД в production
5. Регулярно обновляйте Docker образы

## Резервное копирование

### Бэкап базы данных

```bash
# Создать бэкап
docker compose exec postgres pg_dump -U postgres dashboard > backup.sql

# Восстановить из бэкапа
docker compose exec -T postgres psql -U postgres dashboard < backup.sql
```

## Мониторинг

```bash
# Использование ресурсов
docker stats

# Статус контейнеров
docker compose ps

# Логи в реальном времени
docker compose logs -f
```

