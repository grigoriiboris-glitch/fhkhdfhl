-- Добавляем новые колонки в таблицу posts
ALTER TABLE posts ADD COLUMN IF NOT EXISTS user_id INTEGER;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE posts ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Создаем индекс для быстрого поиска по user_id
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);

-- Создаем индекс для сортировки по времени создания
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at); 