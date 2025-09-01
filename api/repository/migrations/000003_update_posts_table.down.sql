-- Удаляем добавленные колонки
ALTER TABLE posts DROP COLUMN IF EXISTS user_id;
ALTER TABLE posts DROP COLUMN IF EXISTS created_at;
ALTER TABLE posts DROP COLUMN IF EXISTS updated_at;

-- Удаляем индексы
DROP INDEX IF EXISTS idx_posts_user_id;
DROP INDEX IF EXISTS idx_posts_created_at; 