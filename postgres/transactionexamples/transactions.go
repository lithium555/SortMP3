package transactionexamples

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
}

// Create creates new record in DB
func Create(ctx context.Context, dbc *sql.DB) (int64, error) {
	tx, err := dbc.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

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
