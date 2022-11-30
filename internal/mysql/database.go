package mysql

import (
	"context"
	"database/sql"
	"fmt"
)

func DetectDatabase(ctx context.Context, db *sql.DB, dbName string) error {
	rows, err := db.QueryContext(ctx, fmt.Sprintf(DetectDatabaseTemplate, dbName))
	if err != nil {
		return err
	}

	var result string
	for rows.Next() {
		if err := rows.Scan(&result); err != nil {
			return err
		}
		if result == dbName {
			return fmt.Errorf("schema is already exists")
		}
		return nil
	}
	return nil
}
