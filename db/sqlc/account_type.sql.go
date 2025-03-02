// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: account_type.sql

package db

import (
	"context"
)

const createAccountType = `-- name: CreateAccountType :one
INSERT INTO account_types (
  account_type,
  description
) VALUES (
  $1, $2
) RETURNING account_type, description, is_active, created_at, updated_at
`

type CreateAccountTypeParams struct {
	AccountType string `json:"account_type"`
	Description string `json:"description"`
}

func (q *Queries) CreateAccountType(ctx context.Context, arg CreateAccountTypeParams) (AccountType, error) {
	row := q.db.QueryRow(ctx, createAccountType, arg.AccountType, arg.Description)
	var i AccountType
	err := row.Scan(
		&i.AccountType,
		&i.Description,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAccountType = `-- name: DeleteAccountType :exec
UPDATE account_types 
SET is_active = false
WHERE account_type = $1
`

func (q *Queries) DeleteAccountType(ctx context.Context, accountType string) error {
	_, err := q.db.Exec(ctx, deleteAccountType, accountType)
	return err
}

const getAccountType = `-- name: GetAccountType :one
SELECT account_type, description, is_active, created_at, updated_at FROM account_types 
WHERE account_type = $1 AND is_active = true
`

func (q *Queries) GetAccountType(ctx context.Context, accountType string) (AccountType, error) {
	row := q.db.QueryRow(ctx, getAccountType, accountType)
	var i AccountType
	err := row.Scan(
		&i.AccountType,
		&i.Description,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const hardDeleteAccountType = `-- name: HardDeleteAccountType :exec
DELETE FROM account_types
WHERE account_type = $1
`

func (q *Queries) HardDeleteAccountType(ctx context.Context, accountType string) error {
	_, err := q.db.Exec(ctx, hardDeleteAccountType, accountType)
	return err
}

const listAccountTypes = `-- name: ListAccountTypes :many
SELECT account_type, description, is_active, created_at, updated_at FROM account_types
WHERE is_active = true
`

func (q *Queries) ListAccountTypes(ctx context.Context) ([]AccountType, error) {
	rows, err := q.db.Query(ctx, listAccountTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AccountType{}
	for rows.Next() {
		var i AccountType
		if err := rows.Scan(
			&i.AccountType,
			&i.Description,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccountType = `-- name: UpdateAccountType :one
UPDATE account_types
SET 
    account_type = COALESCE($1, account_type),
    description = COALESCE($2, description),
    is_active = COALESCE($3, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE account_type = $4
RETURNING account_type, description, is_active, created_at, updated_at
`

type UpdateAccountTypeParams struct {
	AccountType   string `json:"account_type"`
	Description   string `json:"description"`
	IsActive      bool   `json:"is_active"`
	AccountType_2 string `json:"account_type_2"`
}

func (q *Queries) UpdateAccountType(ctx context.Context, arg UpdateAccountTypeParams) (AccountType, error) {
	row := q.db.QueryRow(ctx, updateAccountType,
		arg.AccountType,
		arg.Description,
		arg.IsActive,
		arg.AccountType_2,
	)
	var i AccountType
	err := row.Scan(
		&i.AccountType,
		&i.Description,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
