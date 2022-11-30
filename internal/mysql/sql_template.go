package mysql

var (
	DetectUserTemplate string = "SELECT EXISTS(SELECT 1 FROM mysql.user WHERE user = '%s');"

	DetectDatabaseTemplate string = "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = '%s';"

	CreateUserTemplate = `CREATE USER IF NOT EXISTS '%s'@'%s' IDENTIFIED BY '%s';`

	CreateDatabaseTemplate = `CREATE DATABASE %s CHARACTER SET %s COLLATE %s`

	GrantUserTemplate = "GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%s' WITH GRANT OPTION;"

	//CreateDBaasManagerTemplate = `CREATE TABLE dbaas.manager (
	//								uid varchar(100) NOT NULL,
	//								name varchar(100) NOT NULL,
	//								db_name varchar(100) NOT NULL,
	//								type varchar(100) NOT NULL,
	//								username varchar(100) NOT NULL,
	//								password varchar(100) NOT NULL,
	//								ip_range varchar(100) NOT NULL,
	//								charset varchar(100) NOT NULL,
	//								collation varchar(100) NOT NULL,
	//								PRIMARY KEY (uid)
	//                       		)
	//							ENGINE=InnoDB
	//							DEFAULT CHARSET=utf8mb4
	//							COLLATE=utf8mb4_general_ci;`
	//

	//
	//UpdatePasswordTemplate = "ALTER USER '%s'@'%s' IDENTIFIED BY '%s';"
	//
	//DeleteUserTemplate = `DELETE FROM mysql.user WHERE Host='%s' AND User='%s';`
	//
	//DetectUserTemplate string = "SELECT EXISTS(SELECT 1 FROM mysql.user WHERE user = '%s');"
)
