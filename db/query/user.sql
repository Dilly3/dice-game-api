-- name: CreateUser :one
INSERT INTO users (
  firstname, lastname, username, password
) VALUES (
  $1, $2 , $3 , $4 
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

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1
LIMIT 1;


-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;

-- name: UpdateUserGameMode :exec
UPDATE users 
 set game_mode = $2  
WHERE username = $1;



