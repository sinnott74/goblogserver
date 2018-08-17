package orm

import (
	"strconv"
	"strings"
)

// deleteByIDQuery creates a Delete query which delete by a rows ID field
// e.g. DELETE FROM test WHERE id=$1
func deleteByIDQuery(tableName string) string {
	return deleteQuery(tableName, []string{"id"})
}

// selectQuery creates a select db query
// i.e. SELECT col1, col2, col3 FROM table WHERE col1=$1 AND col3=$2;
// or SELECT * FROM table; if no columneNames or whereNames are given
func getByIDQuery(tableName string) string {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT * FROM ")
	queryBuilder.WriteString(tableName)
	queryBuilder.WriteString(" WHERE id=$1 LIMIT 1;")
	return queryBuilder.String()
}

// selectQuery creates a select db query
// i.e. SELECT col1, col2, col3 FROM table WHERE col1=$1 AND col3=$2;
// or SELECT * FROM table; if no columneNames or whereNames are given
func selectQuery(tableName string, columnNames []string, whereNames []string) string {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT ")
	if columnNames == nil || len(columnNames) == 0 {
		queryBuilder.WriteString("*")
	} else {
		for i, columnName := range columnNames {
			if i != 0 {
				queryBuilder.WriteString(", ")
			}
			queryBuilder.WriteString(columnName)
		}
	}
	queryBuilder.WriteString(" FROM ")
	queryBuilder.WriteString(tableName)
	if len(whereNames) == 0 {
		queryBuilder.WriteString(";")
		return queryBuilder.String()
	}
	queryBuilder.WriteString(" WHERE ")
	for i, whereName := range whereNames {
		if i != 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(whereName)
		queryBuilder.WriteString("=$")
		queryBuilder.WriteString(strconv.Itoa(i + 1))
	}
	queryBuilder.WriteString(";")
	return queryBuilder.String()
}

// insertQuery creates the INSERT query string
// i.e INSERT INTO (col1, col2, col3) VALUES (val1, val2, val3)
func insertQuery(tableName string, columnNames []string) string {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO ")
	queryBuilder.WriteString(tableName)
	queryBuilder.WriteString(" (")
	for i, columnName := range columnNames {
		if i != 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString(columnName)
	}
	queryBuilder.WriteString(") VALUES (")
	for i := 0; i < len(columnNames); i++ {
		if i != 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString("$")
		queryBuilder.WriteString(strconv.Itoa(i + 1))
	}
	queryBuilder.WriteString(") RETURNING id;")
	return queryBuilder.String()
}

// deleteQuery creates a query to delete a row from a database table
// e.g. DELETE FROM test WHERE col1=$1 AND col3=$2;
func deleteQuery(tableName string, whereNames []string) string {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("DELETE FROM ")
	queryBuilder.WriteString(tableName)
	if len(whereNames) == 0 {
		queryBuilder.WriteString(";")
		return queryBuilder.String()
	}
	queryBuilder.WriteString(" WHERE ")
	for i, whereName := range whereNames {
		if i != 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(whereName)
		queryBuilder.WriteString("=$")
		queryBuilder.WriteString(strconv.Itoa(i + 1))
	}
	queryBuilder.WriteString(";")
	return queryBuilder.String()
}

// updateQuery create an update sql query
// e.g. UPDATE tableName SET col1=$1, col2=$2; or UPDATE tableName SET col1=$1, col2=$2 WHERE col1=$3;
func updateQuery(tableName string, setNames []string, whereNames []string) string {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("UPDATE ")
	queryBuilder.WriteString(tableName)
	queryBuilder.WriteString(" SET ")
	for i, setName := range setNames {
		if i != 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString(setName)
		queryBuilder.WriteString("=$")
		queryBuilder.WriteString(strconv.Itoa(i + 1))
	}
	if len(whereNames) == 0 {
		queryBuilder.WriteString(";")
		return queryBuilder.String()
	}
	queryBuilder.WriteString(" WHERE ")
	varOffset := 1 + len(setNames)
	for i, whereName := range whereNames {
		if i != 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(whereName)
		queryBuilder.WriteString("=$")
		queryBuilder.WriteString(strconv.Itoa(i + varOffset))
	}
	queryBuilder.WriteString(";")
	return queryBuilder.String()
}

// countQuery creates an SQL query to count
// i.e. SELECT COUNT(*) FROM tablename; or SELECT COUNT(*) FROM tablename WHERE col1=1$ AND col2=$2;
func countQuery(tableName string, whereNames []string) string {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT COUNT(*) as count FROM ")
	queryBuilder.WriteString(tableName)
	if len(whereNames) == 0 {
		queryBuilder.WriteString(";")
		return queryBuilder.String()
	}
	queryBuilder.WriteString(" WHERE ")
	for i, whereName := range whereNames {
		if i != 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(whereName)
		queryBuilder.WriteString("=$")
		queryBuilder.WriteString(strconv.Itoa(i + 1))
	}
	queryBuilder.WriteString(";")
	return queryBuilder.String()
}
