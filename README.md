
# API организационной структуры

## Быстрый старт

1. **Клонируй репозиторий**.
   
2. **Создай `.env` из примера**.  

3. **Запусти сервисы**:
   ```bash
   docker-compose up --build
   ```
   Будут запущены PostgreSQL, применены миграции и стартует API на порту `8080`.

4. **Проверь работу**, отправив тестовый запрос:
   ```bash
   curl -i -X POST http://localhost:8080/departments \
     -H "Content-Type: application/json" \
     -d '{"name":"IT"}'
   ```
## API Эндпоинты

### Создать подразделение
```http
POST /departments
Content-Type: application/json

{
  "name": "IT",
  "parent_id": null   
}
```
Ответ: `201 Created` (объект подразделения)

### Создать сотрудника в подразделении
```http
POST /departments/{id}/employees
Content-Type: application/json

{
  "full_name": "John Doe",
  "position": "Developer",
  "hired_at": "2025-01-15"   
}
```
Ответ: `201 Created` (объект сотрудника)

### Получить подразделение с деревом и сотрудниками
```http
GET /departments/{id}?depth=2&include_employees=true
```
Ответ: `200 OK` (объект подразделения с полями `children` и `employees`)

### Переместить / переименовать подразделение
```http
PATCH /departments/{id}
Content-Type: application/json

{
  "name": "Новое имя",        
  "parent_id": 2              
}
```
Ответ: `200 OK` (обновлённый объект)

### Удалить подразделение
```http
DELETE /departments/{id}?mode=cascade
```
Ответ: `204 No Content`


## Архитектура
```
cmd/           # Точка входа
internal/
  config             # Загрузка конфигурации
  model              # Модели GORM (Department, Employee)
  repository         # Доступ к БД (интерфейс + реализация)
  service            # Бизнес-логика и валидация
  handler            # HTTP-обработчики
  router             # Настройка маршрутов
  middleware          # Логирование
db/migrate           # SQL-миграции goose
```


