BEGIN;

CREATE TABLE IF NOT EXISTS refresh_tokens (
    user_id BIGINT NOT NULL,
    ip VARCHAR(255) NOT NULL,
    token VARCHAR(60) NOT NULL,
    expiration_time TIMESTAMP NOT NULL,
    PRIMARY KEY(user_id, ip)
);

COMMIT;
