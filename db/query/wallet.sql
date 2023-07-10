
-- name: CreateWallet :one
INSERT INTO wallets (
  user_id, username
) VALUES (
  $1, $2 
)
RETURNING *;

-- name: DeleteWallet :exec
DELETE FROM wallets
WHERE username = $1;

-- name: GetWalletByUsername :one
SELECT * FROM wallets
WHERE username = $1;

-- name: GetWalletByUsernameForUpdate :one
SELECT * FROM wallets
WHERE username = $1
FOR UPDATE;

-- name: UpdateWallet :exec
UPDATE wallets
  set balance = $2
WHERE username = $1;
