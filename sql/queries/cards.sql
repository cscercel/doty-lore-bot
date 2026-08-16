-- name: CreateCard :one
INSERT INTO lore_cards (name, type, summary, body, image_url, tags)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetCardByName :one
SELECT * 
FROM lore_cards 
WHERE name = ? collate nocase;

-- name: ListCards :many
SELECT *
FROM lore_cards
ORDER BY name;

-- name: ListCardsByType :many
SELECT *
FROM lore_cards
WHERE type = ?
ORDER BY name;

-- name: UpdateCard :one
UPDATE lore_cards
SET 
    name = ?,
    summary = ?,
    body = ?,
    tags = ?,
    image_url = ?,
    updated_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: UpdateCardImage :one
UPDATE lore_cards
SET 
    image_url = ?,
    updated_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: SearchCardsByName :many
SELECT id, name, type
FROM lore_cards
WHERE name LIKE ? || '%'
ORDER BY name
LIMIT 25;

-- name: DeleteCard :exec
DELETE FROM lore_cards
WHERE id = ?;

-- name: DeleteCardByName :exec
DELETE FROM lore_cards
WHERE name = ? collate nocase;
