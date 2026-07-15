CREATE TABLE habit_logs (
    id BIGSERIAL PRIMARY KEY,

    habit_id BIGINT NOT NULL,

    completed_at TIMESTAMP NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT habit_logs_habit_id_fkey
        FOREIGN KEY (habit_id)
        REFERENCES habits(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_habit_logs_unique_completion
ON habit_logs (
    habit_id,
    ((completed_at)::date)
);