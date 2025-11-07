-- name: CreateSummary :one
INSERT INTO summaries (
  user_username, recommendation_id, pdf_path
) VALUES ($1, $2, $3)
RETURNING *;

-- name: ListSummaries :many
SELECT * FROM summaries
WHERE user_username = $1
ORDER BY id DESC;

-- name: GetSummary :one
SELECT * FROM summaries WHERE id = $1 LIMIT 1;

-- name: DeleteSummary :exec
DELETE FROM summaries WHERE id = $1;
