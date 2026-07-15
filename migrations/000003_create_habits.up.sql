CREATE TABLE habits (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,

    title VARCHAR(255) NOT NULL,

    description TEXT,

    color VARCHAR(50),

    icon VARCHAR(100),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMP,

    CONSTRAINT habits_user_id_fkey
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);