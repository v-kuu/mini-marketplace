package sqlite
import (
	"context"
	"database/sql"
)

func withTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:all

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
