-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id, amount, transaction_type, username
) VALUES (
  $1, $2 , $3 , $4
)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE user_id = $1
AND transaction_type = $2;

-- name: GetTransactionsByUsername :many
SELECT * FROM transactions
WHERE username = $1;