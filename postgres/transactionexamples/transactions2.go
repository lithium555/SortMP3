package transactionexamples

import (
	"context"
	"database/sql"
	"errors"
)

// Create2 creates new record in DB
func Create2(ctx context.Context, dbc *sql.DB) (int64, error) {
	tx, err := dbc.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // nolint yes we can do not check error, when we use defer

	completed := 10
	price := 100

	res, err := tx.ExecContext(ctx, `some sql insert query`,
		completed, price)

	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	} else if count == 0 {
		return 0, errors.New("no rows updated")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

// SomeType represents int64 type
type SomeType struct {
	ID int64
}

// Lookup searches Point element in DB by id
func Lookup(ctx context.Context, dbc *sql.DB, id int64) (*SomeType, error) {
	var st SomeType

	row := dbc.QueryRowContext(ctx, `some select query`)

	if err := row.Scan(&st.ID); err != nil {
		return nil, err
	}

	return &st, nil
}
