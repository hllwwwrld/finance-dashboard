# Finance Dashboard — развёртывание

Документ описывает два способа запуска проекта и как они устроены.


| Режим                    | Когда использовать                    | Команды                   |
| ------------------------ | ------------------------------------- | ------------------------- |
| **Локальная разработка** | Разработка с hot-reload               | `make dev-`*              |
| **Production (Docker)**  | Сервер, staging, первый полный запуск | `make -f Makefile.prod `* |


Рабочая директория для всех команд — **корень репозитория** (`finance-dashboard/`).

---

## Архитектура

### Локальная разработка

Go-бэкенд и Next.js-фронтенд запускаются **на хост-машине**. В Docker поднимается только PostgreSQL.

```
┌─────────────────────────────────────────────────────────┐
│  Хост (ваша машина)                                     │
│                                                         │
│  Browser ──► :3001  Next.js (pnpm dev)                  │
│                    │                                    │
│                    │ NEXT_PUBLIC_API_URL=               │
│                    │ http://localhost:3000              │
│                    ▼                                    │
│              :3000  Go backend (go run)                 │
│                    │                                    │
│                    ▼                                    │
│              :5432  PostgreSQL (Docker)                 │
└─────────────────────────────────────────────────────────┘
```

**Особенности:**

- UI доступен на `http://localhost:3001`
- API — на `http://localhost:3000` (фронт ходит напрямую, без nginx)
- CORS разрешает `localhost` и private IP автоматически
- `.env` в корне **не подхватывается** `go run` автоматически — backend использует defaults (`POSTGRES_HOST=localhost`)



### Production (Docker Compose)

Полный стек в контейнерах. Единственная точка входа снаружи — **nginx**.

```
                    ┌────────── nginx ──────────┐
 Browser ──► :8080  │  /api/*  → backend:3000   │
 (или :443 HTTPS)   │  /*      → frontend:3000  │
                    └───────────────────────────┘
                              │
         ┌────────────────────┼────────────────────┐
         ▼                    ▼                    ▼
     backend:3000       frontend:3000         postgres:5432
         │                                         ▲
         └─────────────────────────────────────────┘
                              ▲
                         migrations (goose, one-shot)
```

**Порядок запуска:** `postgres` → `migrations` → `backend` → `frontend` → `nginx`

**Маршрутизация nginx:**

- `GET /api/ping` → backend `GET /ping` (префикс `/api` срезается)
- `GET /` → Next.js production server

Frontend собран с `NEXT_PUBLIC_API_URL=/api` — браузер обращается к тому же origin через nginx.

---



## Требования



### Локальная разработка

- Docker (только для Postgres)
- Go 1.25+
- Node.js 20+, pnpm
- [goose](https://github.com/pressly/goose) (`go install github.com/pressly/goose/v3/cmd/goose@latest`)



### Production

- Docker Engine + Docker Compose (v2)
- Make (опционально, для удобных команд)

Проверка:

```bash
docker --version
docker compose version
```

---



## Переменные окружения

Скопируйте шаблон:

```bash
cp env.example .env
```


| Переменная          | Локально (`go run`)   | Docker (prod)                            | Описание                                      |
| ------------------- | --------------------- | ---------------------------------------- | --------------------------------------------- |
| `POSTGRES_HOST`     | `localhost` (default) | `postgres`                               | Хост БД                                       |
| `POSTGRES_PORT`     | `5432`                | `5432`                                   | Порт БД                                       |
| `POSTGRES_USER`     | `postgres`            | `postgres`                               | Пользователь БД                               |
| `POSTGRES_PASSWORD` | `postgres`            | **смените!**                             | Пароль БД                                     |
| `POSTGRES_DB`       | `dashboard`           | `dashboard`                              | Имя БД                                        |
| `DB_SSLMODE`        | `disable`             | `require` для внешней БД                 | SSL-подключение к Postgres                    |
| `JWT_SECRET`        | обязательно           | обязательно                              | Секрет для JWT (не используйте default `123`) |
| `SECURE_COOKIE`     | `false`               | `true` при HTTPS                         | Флаг `Secure` для auth-cookie                 |
| `NGINX_PORT`        | —                     | `8080` (default)                         | HTTP-порт nginx на хосте                      |
| `NGINX_HTTPS_PORT`  | —                     | `8443` (default, на сервере часто `443`) | HTTPS-порт nginx                              |


> **Важно:** в `env.example` для Docker указан `POSTGRES_HOST=postgres`. Для локального `go run` либо не задавайте эту переменную (сработает default `localhost`), либо явно установите `POSTGRES_HOST=localhost`.

---



## Локальная разработка



### Быстрый старт (всё сразу)

```bash
make dev-all
```

Поднимет Postgres + миграции, backend и frontend в одном терминале.

### Пошагово (отдельные терминалы)

```bash
# Терминал 1 — БД и миграции
make dev-db

# Терминал 2 — backend
make dev-backend

# Терминал 3 — frontend
make dev-frontend
```



### Доступ


| Сервис   | URL                                                      |
| -------- | -------------------------------------------------------- |
| UI       | [http://localhost:3001](http://localhost:3001)           |
| API      | [http://localhost:3000/ping](http://localhost:3000/ping) |
| Postgres | localhost:5432                                           |




### Доступ с телефона в локальной сети

```bash
make dev-frontend-lan LAN_IP=192.168.x.x
# backend слушает 0.0.0.0:3000 (go run делает это по умолчанию)
```



### Передача env в локальный backend

`make dev-backend` не загружает `.env`. Варианты:

```bash
# разово
export $(grep -v '^#' .env | xargs) && make dev-backend

# или только нужные переменные
POSTGRES_HOST=localhost JWT_SECRET=your-secret make dev-backend
```

---



## Production (Docker)



### Первый запуск

```bash
# 1. Настройка окружения
cp env.example .env
# Отредактируйте: JWT_SECRET, POSTGRES_PASSWORD, при HTTPS — SECURE_COOKIE=true

# 2. Сборка образов
make -f Makefile.prod build

# 3. Запуск
make -f Makefile.prod up
```



### Проверка

```bash
docker compose ps
docker compose logs migrations
docker compose logs backend

# API через nginx
curl http://localhost:8080/api/ping

# UI
open http://localhost:8080
```



### HTTPS

**Самоподписанный сертификат (dev/staging):**

```bash
make -f Makefile.prod ssl-selfsigned SSL_DOMAIN=your-domain.com
make -f Makefile.prod up-https
# https://localhost:8443
```

**Let's Encrypt / свой сертификат:**

```bash
# Положите файлы:
#   nginx/ssl/fullchain.pem
#   nginx/ssl/privkey.pem
make -f Makefile.prod up-https
```

На сервере добавьте в `.env`:

```
NGINX_HTTPS_PORT=443
SECURE_COOKIE=true
```



### Basic Auth (опционально)

По умолчанию отключён (закомментирован в `nginx/snippets/app-locations.conf`).

Для включения:

1. Создайте пароль: `make -f Makefile.prod nginx-htpasswd HTPASSWD_USER=admin HTPASSWD_PASS=secret`
2. Раскомментируйте `auth_basic` в `nginx/snippets/app-locations.conf`
3. Пересоберите nginx: `make -f Makefile.prod build-nc && make -f Makefile.prod up`



### Управление

```bash
make -f Makefile.prod down       # остановить
make -f Makefile.prod down-v     # остановить + удалить volume БД
make -f Makefile.prod build-nc   # пересборка без кэша
make -f Makefile.prod help       # все команды
```

---



## Production: внешняя база данных

Если Postgres managed (RDS, Supabase и т.д.), обновите `.env`:

```bash
POSTGRES_HOST=your-db-host.example.com
POSTGRES_PORT=5432
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=your_database
DB_SSLMODE=require
```

Уберите сервис `postgres` из `docker-compose.yml` или создайте override-файл `docker-compose.prod.yml`, где `backend` и `migrations` не зависят от локального `postgres`.

---



## Автозапуск на сервере (systemd)

```ini
# /etc/systemd/system/finance-dashboard.service
[Unit]
Description=Finance Dashboard
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/path/to/finance-dashboard
ExecStart=/usr/bin/make -f Makefile.prod up-https
ExecStop=/usr/bin/make -f Makefile.prod down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable finance-dashboard
sudo systemctl start finance-dashboard
```

---



## Резервное копирование БД

```bash
# Создать бэкап
docker compose exec postgres pg_dump -U postgres dashboard > backup.sql

# Восстановить
docker compose exec -T postgres psql -U postgres dashboard < backup.sql
```

---



## Устранение проблем



### Backend не подключается к БД (Docker)

```bash
docker compose logs backend
docker compose logs postgres
```

Проверьте, что в `.env` для Docker: `POSTGRES_HOST=postgres` (не `localhost`).

### Backend не подключается к БД (локально)

- Postgres должен быть запущен: `make dev-db`
- `POSTGRES_HOST` должен быть `localhost`, не `postgres`



### Порт 5432 занят

```bash
# В .env
POSTGRES_PORT=5433
```



### Порт 8080 занят

```bash
# В .env
NGINX_PORT=8081
```



### Миграции не применились

```bash
docker compose logs migrations
# Перезапуск migrations вручную:
docker compose run --rm migrations
```



### Auth не работает в Docker

Убедитесь, что `JWT_SECRET` задан в `.env` (не оставляйте default `123` в production).

### Frontend бьёт не в тот API


| Режим       | `NEXT_PUBLIC_API_URL`        |
| ----------- | ---------------------------- |
| Docker prod | `/api` (задаётся при сборке) |
| Local dev   | `http://localhost:3000`      |


Не используйте `NEXT_PUBLIC_API_URL=/api` локально без nginx — запросы пойдут в mock-routes Next.js (`frontend/app/api/`).

---



## Справочник команд

```bash
make help                          # локальная разработка
make -f Makefile.prod help         # production Docker

# Локально
make dev-db | dev-backend | dev-frontend | dev-all

# Production
make -f Makefile.prod build | up | up-https | down | down-v
make -f Makefile.prod ssl-selfsigned
make -f Makefile.prod nginx-htpasswd HTPASSWD_USER=u HTPASSWD_PASS=p
```

---



## Безопасность (production checklist)

- [ ] Сменить `POSTGRES_PASSWORD` и `JWT_SECRET`
- [ ] `DB_SSLMODE=require` для внешней БД
- [ ] `SECURE_COOKIE=true` при HTTPS
- [ ] Настроить HTTPS (Let's Encrypt)
- [ ] Не коммитить `.env`, `nginx/.htpasswd`, `nginx/ssl/*.pem`
- [ ] Ограничить доступ к порту Postgres (`POSTGRES_PORT`) firewall'ом
- [ ] Рассмотреть включение Basic Auth для staging