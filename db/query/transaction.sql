-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id, amount, transaction_type
) VALUES (
  $1, $2 , $3
)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE user_id = $1
AND transaction_type = $2;

-- name: GetTransactionByUserId :one
SELECT * FROM transactions
WHERE user_id = $1;