-- name: GetWalletByID :one
SELECT * FROM wallets WHERE id = $1;

-- name: UpdateWallet :one
UPDATE wallets 
SET balance = balance + @amount
WHERE id = @id
RETURNING *;
