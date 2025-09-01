CREATE TABLE IF NOT EXISTS mindmaps (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    data TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_mindmaps_user_id ON mindmaps(user_id);
CREATE INDEX idx_mindmaps_is_public ON mindmaps(is_public);
CREATE INDEX idx_mindmaps_updated_at ON mindmaps(updated_at); 