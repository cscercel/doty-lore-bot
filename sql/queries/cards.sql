-- name: CreateCard :one
INSERT INTO lore_cards (name, type, summary, body, image_url, tags)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCardByName :one
SELECT * 
FROM lore_cards 
WHERE LOWER(name) = LOWER($1);

-- name: ListCards :many
SELECT *
FROM lore_cards
ORDER BY name;

-- name: ListCardsByType :many
SELECT *
FROM lore_cards
WHERE type = $1
ORDER BY name;

-- name: UpdateCard :one
UPDATE lore_cards
SET 
    summary = $2,
    body = $3,
    tags = $4,
    image_url = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateCardImage :one
UPDATE lore_cards
SET 
    image_url = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SearchCardsByName :many
SELECT id, name, type
FROM lore_cards
WHERE name ILIKE $1 || '%'
ORDER BY name
LIMIT 25;

-- name: DeleteCard :exec
DELETE FROM lore_cards
WHERE id = $1;

-- name: DeleteCardByName :exec
DELETE FROM lore_cards
WHERE LOWER(name) = LOWER($1);
