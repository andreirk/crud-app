DROP TABLE IF EXISTS refresh_tokens;

CREATE TABLE refresh_tokens (
    id SERIAL NOT NULL UNIQUE,
    user_id INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL
);
