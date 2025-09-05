# MindMap с авторизацией

Веб-приложение для создания ментальных карт с системой авторизации пользователей.

## Возможности

- ✅ Создание и редактирование ментальных карт
- ✅ Система регистрации и авторизации пользователей
- ✅ Ролевая модель доступа (RBAC)
- ✅ JWT токены для аутентификации
- ✅ Многоязычная поддержка (русский, английский, китайский)
- ✅ PostgreSQL для хранения данных
- ✅ Docker Compose для развертывания

## Технологии

### Frontend
- Vue.js 3
- Vue Router
- Vue I18n
- Simple Mind Map

### Backend
- Go
- PostgreSQL
- Casbin (RBAC)
- JWT

## Быстрый запуск

### 1. Клонирование репозитория
```bash
git clone <repository-url>
cd hyy-vue3-mindMap
```

### 2. Запуск с Docker Compose
```bash
# Запуск всех сервисов
docker-compose up -d

# Просмотр логов
docker-compose logs -f

# Остановка
docker-compose down
```

### 3. Доступ к приложению
- Веб-интерфейс: http://localhost:8080
- База данных: localhost:5432

## Разработка

### Запуск в режиме разработки
```bash
# Установка зависимостей
npm install

# Запуск только бэкенда
cd api
docker-compose up -d postgres
go run main.go

# Запуск фронтенда (в другом терминале)
npm run serve
```

### Структура проекта
```
hyy-vue3-mindMap/
├── src/
│   ├── components/
│   │   ├── Auth/           # Компоненты авторизации
│   │   │   ├── Login.vue
│   │   │   ├── Register.vue
│   │   │   └── AuthWrapper.vue
│   │   └── MindMapList.vue # Список ментальных карт
│   ├── lang/               # Файлы локализации
│   │   ├── ru_ru.js        # Русский
│   │   ├── en_us.js        # Английский
│   │   └── zh_cn.js        # Китайский
│   ├── utils/
│   │   └── auth.js         # Утилиты авторизации
│   ├── pages/              # Страницы приложения
│   └── App.vue             # Главный компонент
├── api/                    # Go бэкенд
│   ├── auth/              # Сервис авторизации
│   ├── models/            # Модели данных
│   │   ├── user.go
│   │   ├── post.go
│   │   └── mindmap.go
│   ├── repository/        # Слой доступа к данным
│   │   ├── migrations/    # Миграции БД
│   │   ├── user_repository.go
│   │   ├── post_repository.go
│   │   └── mindmap_repository.go
│   └── main.go           # Главный файл
├── docker-compose.yaml    # Конфигурация Docker
├── nginx.conf            # Конфигурация nginx
└── Dockerfile            # Docker образ
```

## API Endpoints

### Авторизация
- `POST /auth/login` - Вход в систему
- `POST /auth/register` - Регистрация
- `GET /auth/logout` - Выход из системы
- `GET /auth/check` - Проверка авторизации
- `GET /auth/user` - Информация о пользователе

### Ментальные карты
- `GET /api/mindmaps` - Список карт пользователя
- `POST /api/mindmaps` - Создание новой карты
- `GET /api/mindmaps/{id}` - Получение карты
- `PUT /api/mindmaps/{id}` - Обновление карты
- `DELETE /api/mindmaps/{id}` - Удаление карты
- `GET /api/mindmaps/public` - Список публичных карт

## Роли и права доступа

### Пользователь (user)
- Создание и редактирование своих ментальных карт
- Просмотр публичных карт

### Автор (author)
- Все права пользователя
- Публикация карт

### Администратор (admin)
- Все права
- Управление пользователями
- Доступ ко всем картам

## База данных

### Таблица users
- `id` - Уникальный идентификатор
- `name` - Имя пользователя
- `email` - Email (уникальный)
- `password` - Хешированный пароль
- `role` - Роль пользователя
- `created_at` - Время создания
- `updated_at` - Время обновления

### Таблица mindmaps
- `id` - Уникальный идентификатор
- `title` - Заголовок карты
- `data` - JSON данные карты
- `user_id` - ID владельца карты
- `is_public` - Публичная ли карта
- `created_at` - Время создания
- `updated_at` - Время обновления

## Безопасность

- Пароли хешируются с помощью bcrypt
- JWT токены для аутентификации
- RBAC модель доступа с Casbin
- Проверка прав на уровне middleware
- Защита от CSRF атак

## Конфигурация

### Переменные окружения
Создайте файл `.env` в папке `api/`:
```env
POSTGRES_DB=mindmap_db
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_HOST=localhost:5432
PORT=8000
LOG_LEVEL=info
```

### Проксирование в разработке
В `vue.config.js` настроено проксирование запросов к API:
```javascript
devServer: {
  proxy: {
    '/auth': {
      target: 'http://localhost:8000',
      changeOrigin: true
    },
    '/api': {
      target: 'http://localhost:8000',
      changeOrigin: true
    }
  }
}
```

## Развертывание

### Продакшн
```bash
# Сборка и запуск
docker-compose -f docker-compose.yaml up -d

# Проверка статуса
docker-compose ps

# Просмотр логов
docker-compose logs -f mindmap-frontend
docker-compose logs -f api-backend
```

### Миграции базы данных
```bash
# Автоматически выполняются при запуске
# Или вручную:
docker-compose exec api-backend go run main.go
```

## Устранение неполадок

### Проблемы с подключением к БД
```bash
# Проверка статуса PostgreSQL
docker-compose ps postgres

# Просмотр логов
docker-compose logs postgres

# Подключение к БД
docker-compose exec postgres psql -U postgres -d mindmap_db
```

### Проблемы с авторизацией
```bash
# Проверка логов API
docker-compose logs api-backend

# Проверка JWT токенов
# Токены хранятся в cookies
```

### Проблемы с фронтендом
```bash
# Проверка логов nginx
docker-compose logs mindmap-frontend

# Пересборка образа
docker-compose build mindmap-frontend
docker-compose up -d mindmap-frontend
```
