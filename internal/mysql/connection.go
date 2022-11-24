package mysql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	datasourcev1alpha1 "github.com/yuanyp8/dbhero/apis/datasource/v1alpha1"
	"time"
)

const (
	DsnTemplate = "%s:%s@tcp(%s:%d)/?charset=utf8mb4&multiStatements=true"
	PingTimeout = 2 * time.Second
)

func InitPool(req *datasourcev1alpha1.DataSource) (*sql.DB, error) {
	dsn := fmt.Sprintf(DsnTemplate, req.Status.Auth.Username.Value, req.Status.Auth.Password.Value, req.Spec.Connection.MySQL.Access.Host, req.Spec.Connection.MySQL.Access.Port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}

	// initialize the Mysql connection pool
	db.SetMaxOpenConns(req.Spec.Connection.MySQL.PoolConfig.MaxOpenConn)
	db.SetMaxIdleConns(req.Spec.Connection.MySQL.PoolConfig.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(req.Spec.Connection.MySQL.PoolConfig.MaxLifeTime))
	db.SetConnMaxIdleTime(time.Duration(req.Spec.Connection.MySQL.PoolConfig.MaxIdleTime))
	return db, nil
}

// Ping help to check the backend MySQL service status
func Ping(ctx context.Context, db *sql.DB) error {
	ctx2, cancelFunc := context.WithTimeout(ctx, PingTimeout)

	defer cancelFunc()
	return db.PingContext(ctx2)
}
