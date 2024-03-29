# IHATEINSTAGRAM

**IHATEINSTAGRAM** - это проект на Golang для бэкенда, который демонстрирует мои знания и навыки в веб-разработке. Он использует Fiber, быстрый и легкий веб-фреймворк, для обработки HTTP-запросов и ответов. Он использует pgx, драйвер и инструментарий PostgreSQL, для подключения и взаимодействия с базой данных. Он использует сессии для безопасной аутентификации и авторизации, чтобы защитить конечные точки и проверить пользователей. **IHATEINSTAGRAM** разработан с принципами чистой архитектуры. Его можно легко развернуть с помощью Docker, платформы для контейнеризации и оркестрации.

Проект предоставляет функциональность для профилей, публикации контента, загрузки изображений в посты, лайков постов и подписки на других пользователей. **IHATEINSTAGRAM** - это проект, которым я горжусь, и я надеюсь, что он произведет впечатление на потенциальных работодателей и клиентов.

## Как запустить
Прежде чем продолжить, убедитесь, что вы соответствуете следующим требованиям:
- Вы установили последнюю версию Go.
- Вы установили PostgreSQL и создали базу данных для проекта.
- Вы создали файл .env с определенным `POSTGRES_URL`.

Чтобы запустить **IHATEINSTAGRAM**, выполните следующие действия:
1. Клонируйте этот репозиторий: `git clone https://github.com/indetensai/ihateinstagram.git`
2. Перейдите в каталог проекта: `cd ihateinstagram`
3. Установите зависимости: `go mod download`
4. Заполните файл `.env` необходимыми переменными окружения.
5. Соберите исполняемый файл: `go build -o ihateinstagram cmd/ihateinstagram/main.go`

## Как запустить (docker-compose)
Прежде чем продолжить, убедитесь, что вы соответствуете следующим требованиям:
- Вы установили последнюю версию docker(-desktop).

Чтобы запустить **IHATEINSTAGRAM** с помощью docker-compose, выполните следующие действия:
1. Клонируйте этот репозиторий: `git clone https://github.com/indetensai/ihateinstagram.git`
2. Перейдите в каталог проекта: `cd ihateinstagram`
3. Запустите `docker compose up`

## Использование
Чтобы запустить **IHATEINSTAGRAM**, выполните следующие действия:
1. Запустите исполняемый файл: `./ihateinstagram`
2. Сервер будет слушать порт 8080.
3. Чтобы взаимодействовать с чат API, вы можете использовать любой HTTP-клиент на ваш выбор.

**IHATEINSTAGRAM** API имеет следующие конечные точки:
- `POST /user/register`: Создать новую учетную запись пользователя.
- `POST /user/login`: Войти в систему с существующей учетной записью пользователя и получите сессию.
- `DELETE /session/:session_id`: Удалить сессию.
- `POST /user/:user_id/follow`: Подписывает пользователя на пользователя user_id. Требуется аутентификация.
- `POST /user/:user_id/unfollow`: Отменяет подписку на пользователя от пользователя с user_id. Требуется аутентификация.
- `GET /user/:user_id/followers`: Получить подписчиков пользователя user_id.
- `POST /post`: Создать пост. Требуется аутентификация.
- `GET /post/:post_id` Получить пост по post_id. Требуется аутентификация.
- `PATCH /post/:post_id`: Изменить пост post_id. Требуется аутентификация.
- `PUT /post/:post_id/like`: Поставить лайк на пост post_id. Требуется аутентификация.
- `GET /post/:post_id/likes`: Получить лайки на посте post_id.Требуется аутентификация.
- `DELETE /post/:post_id/like`: Удалить лайк на посте post_id. Требуется аутентификация.
- `DELETE /post/:post_id`: Удалить пост post_id. Требуется аутентификация.
- `POST /post/:post_id/image`: Загрузить изображение в пост post_id. Требуется аутентификация.
- `GET /post/:post_id/images`: Получить изображения поста post_id.Требуется аутентификация.
- `GET /post/:post_id/thumbnails`: Получить миниатюры поста post_id.Требуется аутентификация.