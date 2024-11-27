# Music Library API

## Описание
Music Library API — это REST API для управления библиотекой музыкальных произведений. API предоставляет возможность добавлять, удалять и изменять песни, а также получать информацию о них. Кроме того, реализована интеграция с внешним сервисом для получения дополнительных данных о песнях.

## Возможности
1. **Получение данных библиотеки**:
   - Фильтрация по всем полям.
   - Пагинация.

2. **Получение текста песни**:
   - Текст возвращается с пагинацией по куплетам.

3. **CRUD операции с песнями**:
   - Добавление новой песни.
   - Удаление песни.
   - Изменение данных существующей песни.

4. **Интеграция с внешним API**:
   - При добавлении песни происходит запрос к внешнему API для обогащения данных о песне.

5. **Работа с базой данных**:
   - Используется PostgreSQL.
   - Миграции автоматически применяются при старте.

6. **Логирование**:
   - Покрытие кода debug- и info-логами.

7. **Документация**:
   - Автоматически сгенерированный Swagger.

## Технологии
- Язык: **Go**
- База данных: **PostgreSQL**
- Миграции: **go-migrate**
- Документация API: **Swagger/OpenAPI**
- Логирование: **logrus** или аналогичный инструмент

## Установка и запуск

### Предварительные требования
- **Go** версии 1.20 или выше.
- **Docker** и **Docker Compose** (для запуска базы данных).
- **PostgreSQL** установлен локально (если не используется Docker).
- Файл конфигурации `.env`.

##Примеры запросов

###Получение списка песен:

```http
  GET /songs?page=1&limit=10&filter[group]=Muse

##Добавление новой песни:

```http
  POST /songs
  Content-Type: application/json
  
  {
      "group": "Muse",
      "song": "Supermassive Black Hole"
  }

##Удаление песни:

```http
  DELETE /songs/{id}

##Миграции базы данных
Миграции находятся в папке migrations.

##Swagger
Swagger документация автоматически генерируется и доступна по адресу:

```http
  http://localhost:8080/swagger/index.html

##Структура проекта
cmd/ — точка входа для приложения.
config/ — файлы конфигурации.
internal/ — основная бизнес-логика.
migrations/ — миграции базы данных.
docs/ — документация.
