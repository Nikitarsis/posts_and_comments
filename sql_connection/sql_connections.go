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
	tables  map[string]struct{}
}

func GetObject[T any](csql ConnectionSQL, command string, conversion func(*sql.Rows) *T) *T {
	ret, err := csql.getRows(command)
	if err != nil {
		return nil
	}
	return conversion((ret))
}

func CreateConnection(dbname string) (*ConnectionSQL, error) {
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	database, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	ret := func(str string) (*sql.Rows, error) {
		return database.Query(str)
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
		getRows: ret,
		tables:  tableNames,
	}, nil
}
