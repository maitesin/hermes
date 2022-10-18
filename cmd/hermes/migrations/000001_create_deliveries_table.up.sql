CREATE TABLE IF NOT EXISTS deliveries (
    tracking_id TEXT PRIMARY KEY,
    conversation_id BIGINT NOT NULL,
    log TEXT NOT NULL,
    delivered bool DEFAULT false NOT NULL
);