# **Сервис проведения тендеров**
Этот **api-сервис**, который позволит бизнесу создать тендер на оказание каких-либо услуг. А пользователи/другие бизнесы будут предлагать свои выгодные условия для получения данного тендера.

# Для старта
1. Установить переменные окружения
   SERVER_ADDRESS — адрес и порт, который будет слушать HTTP сервер при запуске. Пример: 0.0.0.0:8080.<br>
   POSTGRES_CONN — URL-строка для подключения к PostgreSQL в формате postgres://{username}:{password}@{host}:{5432}/{dbname}.<br>
   POSTGRES_JDBC_URL — JDBC-строка для подключения к PostgreSQL в формате jdbc:postgresql://{host}:{port}/{dbname}.<br>
   POSTGRES_USERNAME — имя пользователя для подключения к PostgreSQL.<br>
   POSTGRES_PASSWORD — пароль для подключения к PostgreSQL.<br>
   POSTGRES_HOST — хост для подключения к PostgreSQL (например, localhost).<br>
   POSTGRES_PORT — порт для подключения к PostgreSQL (например, 5432).<br>
   POSTGRES_DATABASE — имя базы данных PostgreSQL, которую будет использовать приложение.<br>
2. Запустить docker
```bash
docker-compose build
docker-compose up -d
```
3. Поднять миграции (см. инструкцию ниже)

# **Миграции**
Поднять миграции командой с переменными из окружения
```bash
goose postgres "host=${POSTGRES_HOST} user=${POSTGRES_USERNAME} database=${POSTGRES_DATABASE} password=${POSTGRES_PASSWORD}" up
```