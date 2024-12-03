package db_tx

import (
	"context"
	"database/sql"
	"fmt"
)

func CallTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	err = fn(tx)
	if err != nil {
		if errRb := tx.Rollback(); errRb != nil {
			return fmt.Errorf("error on rollback %v, original error %w", errRb, err)
		}
	}

	return tx.Commit()
}
