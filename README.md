# PR Reviewer Assignment Service (Test Task, Fall 2025)

 Для запуска сервиса необходимо ввести команду 
 
 ```bash
 docker compose up -d --build 
 ```
 
 Чтобы изменить переменные окружения стоит создать копию `.env.example` с названием `.env`. Сам файл `.env.example` не удалять!

## Эндпоинты

- POST `/api/v1/team/add`
- GET  `/api/v1/team/get`
- POST `/api/v1/users/setIsActive`
- GET  `/api/v1/users/add`
- POST `/api/v1/pullRequest/create`
- POST `/api/v1/pullRequest/merge`
- POST `/api/v1/pullRequest/reassign`

## Возникшие вопросы/проблемы по заданию

- Были вопросы к идентефикатору users и pull_requests. В API-контракте у них тип данных string и они имеют вид "u1" и "pr-1001". Не до конца понятны ограничения(всегда ли они имеют такой вид или id может быть любой строкой), поэтому валидаторы для входных данных на всех слоях абстракции не были написаны.
