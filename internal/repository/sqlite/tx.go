package sqlite
import (
	"context"
	"database/sql"
	"errors"
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
	if err := tx.Commit(); err != nil {
		if errors.Is(err, sql.ErrTxDone) {
			return nil
		}
		return err
	}
	return nil
}
