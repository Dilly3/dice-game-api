-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id, amount, balance , transaction_type, username
) VALUES (
  $1, $2 , $3 , $4 , $5
)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE user_id = $1
AND transaction_type = $2;

-- name: GetTransactionsByUsername :many
SELECT * FROM transactions
WHERE username = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: DeleteTransactionByUsername :exec
DELETE FROM transactions
WHERE username = $1;

-- name: UpdateTransaction :exec
UPDATE transactions
  set balance = $2 ,
  amount = $3
WHERE username = $1;