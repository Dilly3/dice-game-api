-- name: CreateUser :one
INSERT INTO users (
  firstname, lastname, username, email, password
) VALUES (
  $1, $2 , $3 , $4 , $5
)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;


-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE username = $1
FOR UPDATE;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY lastname ASC
LIMIT $1
OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;


-- name: UpdateUser :exec
UPDATE users
  set balance = $2
WHERE username = $1;
