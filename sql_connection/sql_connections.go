package sqlconnection

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qwerty1234"
)

type ConnectionSQL struct {
	getRows func(string) (*sql.Rows, error)
	exec    func(string) error
	tables  map[string]struct{}
}

func GetObject[T any](csql ConnectionSQL, command string, conversion func(*sql.Rows) T) T {
	ret, err := csql.getRows(command)
	if err != nil {
		return nil
	}
	return conversion((ret))
}

func Exec(csql ConnectionSQL, command string) error {
	return csql.exec(command)
}

func (c ConnectionSQL) HasTable(str string) bool {
	_, check := c.tables[str]
	return check
}

func (c *ConnectionSQL) AddTable(tableName string, tableDescription string) {
	c.tables[tableName] = struct{}{}
	c.exec(fmt.Sprintf("CREATE TABLE %s", tableDescription))
}

func CreateConnection(dbname string) (*ConnectionSQL, error) {
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	database, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	getrws := func(str string) (*sql.Rows, error) {
		return database.Query(str)
	}
	exec := func(str string) error {
		go database.Exec(str)
	}
	rows, errTable := database.Query("SELECT * FROM pg_catalog.pg_tables")
	if errTable != nil {
		return nil, errTable
	}
	tableNames := make(map[string]struct{})
	defer rows.Close()
	for rows.Next() {
		var str string
		errScan := rows.Scan(&str)
		if errScan != nil {
			return nil, errScan
		}
		tableNames[str] = struct{}{}
	}
	return &ConnectionSQL{
		db:      database,
		getRows: getrws,
		exec: ,
		tables:  tableNames,
	}, nil
}
