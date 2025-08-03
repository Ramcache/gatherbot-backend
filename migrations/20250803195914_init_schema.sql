-- +goose Up
CREATE TABLE events (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        type TEXT NOT NULL CHECK (type IN ('meeting', 'kotyol')),
                        title TEXT NOT NULL,
                        description TEXT,
                        date DATE,
                        time TIME,
                        place TEXT,
                        start_month DATE,
                        end_month DATE,
                        amount INT,
                        owner_id BIGINT NOT NULL,
                        max_participants INT,
                        participants BIGINT[] DEFAULT '{}',
                        created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS events;
