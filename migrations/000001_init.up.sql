-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL, -- Логин пользователя
    password_hash TEXT NOT NULL,          -- Хэш пароля
    role VARCHAR(20) NOT NULL,            -- Роль (обычный пользователь или суперпользователь)
    email VARCHAR(100) UNIQUE NOT NULL,   -- Электронная почта
    nickname VARCHAR(50) UNIQUE NOT NULL, -- Никнейм в формате "@никнейм"
    created_at TIMESTAMP DEFAULT NOW()
);

-- Создание таблицы сообщений
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    sender_username VARCHAR(50) NOT NULL,      -- Никнейм отправителя
    recipient_username VARCHAR(50) NOT NULL,  -- Никнейм получателя
    content TEXT,                              -- Текст сообщения
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    CONSTRAINT fk_sender FOREIGN KEY (sender_username) REFERENCES users (username) ON DELETE CASCADE,
    CONSTRAINT fk_recipient FOREIGN KEY (recipient_username) REFERENCES users (username) ON DELETE CASCADE
);

-- Создание таблицы вложений
CREATE TABLE IF NOT EXISTS attachments (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL,                   -- ID сообщения, к которому прикреплено вложение
    file_data BYTEA,                           -- Данные файла
    file_type VARCHAR(50),                     -- Тип файла (например, image/png)
    file_size INT,                             -- Размер файла
    CONSTRAINT fk_message FOREIGN KEY (message_id) REFERENCES messages (id) ON DELETE CASCADE
);

-- Создание таблицы лайков
CREATE TABLE IF NOT EXISTS likes (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL,                   -- ID сообщения
    user_nickname VARCHAR(50) NOT NULL,        -- Никнейм пользователя, который поставил лайк
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_message FOREIGN KEY (message_id) REFERENCES messages (id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_nickname) REFERENCES users (nickname) ON DELETE CASCADE
);

-- Создание таблицы суперлайков
CREATE TABLE IF NOT EXISTS superlikes (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL,                   -- ID сообщения
    user_nickname VARCHAR(50) NOT NULL,        -- Никнейм пользователя, который поставил суперлайк
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_message FOREIGN KEY (message_id) REFERENCES messages (id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_nickname) REFERENCES users (nickname) ON DELETE CASCADE
);
