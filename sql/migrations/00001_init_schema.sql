-- +goose Up
CREATE TABLE lore_cards (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    type        TEXT NOT NULL,
    summary     TEXT NOT NULL,
    body        TEXT,
    image_url   TEXT,
    tags        TEXT[],
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_lore_cards_name ON lore_cards (name);
CREATE INDEX idx_lore_cards_type ON lore_cards (type);
CREATE UNIQUE INDEX idx_lore_cards_name_unique ON lore_cards (LOWER(name));

-- +goose Down
DROP TABLE lore_cards;
