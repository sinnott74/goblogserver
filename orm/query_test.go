package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQuery(t *testing.T) {
	tableName := "test"
	expectedQuery := "SELECT * FROM test WHERE id=$1;"
	query := getQuery(tableName)
	assert.Equal(t, expectedQuery, query, "Get query not as expected")
}

func TestSelectQuery(t *testing.T) {
	tableName := "test"
	columnNames := []string{"col1", "col2", "col3"}
	whereNames := []string{"col1", "col3"}
	expectedQuery := "SELECT col1, col2, col3 FROM test WHERE col1=$1 AND col3=$2;"
	query := selectQuery(tableName, columnNames, whereNames)
	assert.Equal(t, expectedQuery, query, "Select query not as expected")
}

func TestSelectQueryStar(t *testing.T) {
	tableName := "test"
	columnNames := []string{}
	whereNames := []string{"col1", "col3"}
	expectedQuery := "SELECT * FROM test WHERE col1=$1 AND col3=$2;"
	query := selectQuery(tableName, columnNames, whereNames)
	assert.Equal(t, expectedQuery, query, "Select query not as expected")
}

func TestSelectQueryColumnNil(t *testing.T) {
	tableName := "test"
	whereNames := []string{"col1", "col3"}
	expectedQuery := "SELECT * FROM test WHERE col1=$1 AND col3=$2;"
	query := selectQuery(tableName, nil, whereNames)
	assert.Equal(t, expectedQuery, query, "Select query not as expected")
}

func TestSelectQueryWhereEmpty(t *testing.T) {
	tableName := "test"
	columnNames := []string{"col1", "col2", "col3"}
	whereNames := []string{}
	expectedQuery := "SELECT col1, col2, col3 FROM test;"
	query := selectQuery(tableName, columnNames, whereNames)
	assert.Equal(t, expectedQuery, query, "Select query not as expected")
}

func TestSelectQueryWhereNil(t *testing.T) {
	tableName := "test"
	columnNames := []string{"col1", "col2", "col3"}
	expectedQuery := "SELECT col1, col2, col3 FROM test;"
	query := selectQuery(tableName, columnNames, nil)
	assert.Equal(t, expectedQuery, query, "Select query not as expected")
}

func TestSelectQueryStarWhereNil(t *testing.T) {
	tableName := "test"
	expectedQuery := "SELECT * FROM test;"
	query := selectQuery(tableName, nil, nil)
	assert.Equal(t, expectedQuery, query, "Select query not as expected")
}

func TestInsertQuery(t *testing.T) {
	tableName := "test"
	columnNames := []string{"col1", "col2", "col3"}
	expectedQuery := "INSERT INTO test (col1, col2, col3) VALUES ($1, $2, $3) RETURNING id;"
	query := insertQuery(tableName, columnNames)
	assert.Equal(t, expectedQuery, query, "Insert query not as expected")
}

func TestDeleteQuery(t *testing.T) {
	tableName := "test"
	whereNames := []string{"col1", "col3"}
	expectedQuery := "DELETE FROM test WHERE col1=$1 AND col3=$2;"
	query := deleteQuery(tableName, whereNames)
	assert.Equal(t, expectedQuery, query, "Delete query not as expected")
}

func TestDeleteQueryNoWHere(t *testing.T) {
	tableName := "test"
	expectedQuery := "DELETE FROM test;"
	query := deleteQuery(tableName, nil)
	assert.Equal(t, expectedQuery, query, "Delete query not as expected")
}

func TestDeleteByIDQuery(t *testing.T) {
	tableName := "test"
	expectedQuery := "DELETE FROM test WHERE id=$1;"
	query := deleteByIDQuery(tableName)
	assert.Equal(t, expectedQuery, query, "Delete query not as expected")
}

func TestUpdateQuery(t *testing.T) {
	tableName := "test"
	setNames := []string{"col1", "col2", "col3"}
	whereNames := []string{"col1", "col3"}
	expectedQuery := "UPDATE test SET col1=$1, col2=$2, col3=$3 WHERE col1=$4 AND col3=$5;"
	query := updateQuery(tableName, setNames, whereNames)
	assert.Equal(t, expectedQuery, query, "Update query not as expected")
}

func TestCountQuery(t *testing.T) {
	tableName := "test"
	whereNames := []string{"col1", "col3"}
	expectedQuery := "SELECT COUNT(*) as count FROM test WHERE col1=$1 AND col3=$2;"
	query := countQuery(tableName, whereNames)
	assert.Equal(t, expectedQuery, query, "Count with where clause not as expected")
}

func TestCountQueryNoWhere(t *testing.T) {
	tableName := "test"
	expectedQuery := "SELECT COUNT(*) as count FROM test;"
	query := countQuery(tableName, nil)
	assert.Equal(t, expectedQuery, query, "Count without where clause not as expected")
}
