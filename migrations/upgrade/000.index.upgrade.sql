-- tasks table indexes
CREATE INDEX idx_tasks_user_id ON tasks (user_id);
CREATE INDEX idx_tasks_created_at ON tasks (created_at DESC);

-- users table indexes
CREATE INDEX idx_users_created_at ON users (created_at DESC);