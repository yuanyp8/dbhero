package mysql

import (
	"context"
	"database/sql"
	"fmt"
)

func DetectUser(ctx context.Context, db *sql.DB, username string) (value int, err error) {
	rows, err := db.QueryContext(ctx, fmt.Sprintf(DetectUserTemplate, username))
	if err != nil {
		return -1, err
	}

	var result int
	for rows.Next() {
		if err := rows.Scan(&result); err != nil {
			return -1, err
		}
		if result == 0 {
			return 0, fmt.Errorf("user not found")
		}
		return 0, nil
	}

	return -1, fmt.Errorf("unknown error")
}

func CreateUserDatabaseSTMT(ctx context.Context, db *sql.DB, username string, password string, dbName string, character string, collate string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start tx error, %s", err)
	}

	// 通过defer处理事务提交的方式
	// 1. 无报错，则Commit事务
	// 2. 有报错，则Rollback事务
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// create user
	_, err = tx.ExecContext(ctx, fmt.Sprintf(CreateUserTemplate, username, "%", password))
	if err != nil {
		return err
	}

	// create database
	_, err = tx.ExecContext(ctx, fmt.Sprintf(CreateDatabaseTemplate, dbName, character, collate))
	if err != nil {
		return err
	}

	// grant privileges
	_, err = tx.ExecContext(ctx, fmt.Sprintf(GrantUserTemplate, dbName, username, "%"))

	if err != nil {
		return err
	}
	return nil
}
