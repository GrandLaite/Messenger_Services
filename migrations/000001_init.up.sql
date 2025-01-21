CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    recipient_id INT NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS attachments (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL,
    file_data BYTEA,
    file_type VARCHAR(50),
    file_size INT,
    CONSTRAINT fk_message
        FOREIGN KEY(message_id)
        REFERENCES messages(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS likes (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS superlikes (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL,
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);