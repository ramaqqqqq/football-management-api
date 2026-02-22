CREATE TABLE IF NOT EXISTS goals (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGINT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    player_id BIGINT NOT NULL REFERENCES players(id),
    goal_minute INT NOT NULL CHECK (goal_minute >= 1 AND goal_minute <= 120),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_goals_match_id ON goals(match_id);
CREATE INDEX IF NOT EXISTS idx_goals_player_id ON goals(player_id);
CREATE INDEX IF NOT EXISTS idx_goals_deleted_at ON goals(deleted_at);
