CREATE TABLE IF NOT EXISTS players (
    id BIGSERIAL PRIMARY KEY,
    team_id BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    height DECIMAL(5,2) NOT NULL DEFAULT 0,
    weight DECIMAL(5,2) NOT NULL DEFAULT 0,
    position VARCHAR(20) NOT NULL CHECK (position IN ('penyerang', 'gelandang', 'bertahan', 'penjaga_gawang')),
    jersey_number INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_players_team_id ON players(team_id);
CREATE INDEX IF NOT EXISTS idx_players_deleted_at ON players(deleted_at);

CREATE UNIQUE INDEX IF NOT EXISTS idx_players_jersey_team
    ON players(team_id, jersey_number)
    WHERE deleted_at IS NULL;
