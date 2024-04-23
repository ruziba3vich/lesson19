package packages

import (
	"database/sql"
	"fmt"
)

func LoadFromDb(db *sql.DB, dbName string, data *[]string) error {
	query := fmt.Sprintf("SELECT fullname FROM %s;", dbName)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var fullname string
		if err := rows.Scan(&fullname); err != nil {
			return err
		}
		*data = append(*data, fullname)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
