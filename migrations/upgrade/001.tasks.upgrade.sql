CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    user_information JSONB,
    client VARCHAR(255),
    job VARCHAR(255),
    item VARCHAR(255),
    role VARCHAR(100),
    note TEXT,
    ot VARCHAR(255),
    hours VARCHAR(100),
    volume INT,
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    work_time INT DEFAULT 0,
    minute INT DEFAULT 0,
    break_time INT DEFAULT 0,
    status VARCHAR(50),
    delivery BOOLEAN DEFAULT FALSE,
    custom_created_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE tasks
ALTER COLUMN id TYPE BIGINT;

ALTER TABLE tasks
ADD COLUMN hours VARCHAR(100);