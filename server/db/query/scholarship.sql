-- name: CreateScholarship :one
INSERT INTO scholarships (
  user_username, title, description, match_score, link
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: ListScholarshipsByUser :many
SELECT * FROM scholarships
WHERE user_username = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;

-- name: DeleteScholarship :exec
DELETE FROM scholarships
WHERE id = $1 AND user_username = $2;
