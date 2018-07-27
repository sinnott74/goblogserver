package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestEntity struct {
	ID     int
	Name   string
	Column string
}

var testEntity = &TestEntity{1, "Joe", "Age"}

// func TestGet(t *testing.T) {

// 	ctx := context.Background()
// 	transaction := database.NewTransaction(ctx)
// 	ctx = database.SetTransaction(ctx, transaction)

// 	ID := int64(198)
// 	expectedDOB, err := time.Parse(time.RFC3339Nano, "2018-07-20T11:07:48.34352Z")
// 	assert.NoError(t, err, "Error during Get expected DOB parse")
// 	expectedUser := &model.User{ID: ID, Username: "reflect15@test.com", Firstname: "Test", Lastname: "Reflection", DOB: expectedDOB}

// 	user := &model.User{ID: ID}
// 	err = Get(ctx, user)

// 	assert.NoError(t, err, "Error during Get")
// 	assert.Equal(t, expectedUser, user)

// 	transaction.Rollback()
// }

// func TestSelectAll(t *testing.T) {

// 	ctx := context.Background()
// 	transaction := database.NewTransaction(ctx)
// 	ctx = database.SetTransaction(ctx, transaction)

// 	users := &[]model.User{}
// 	where := &model.User{}
// 	err := SelectAll(ctx, users, where)

// 	assert.NoError(t, err, "Error during SelectAll")
// 	assert.True(t, len(*users) > 0, "Users were returned")

// 	transaction.Rollback()
// }

// func TestInsert(t *testing.T) {

// 	// t.SkipNow()

// 	ctx := context.Background()
// 	transaction := database.NewTransaction(ctx)
// 	ctx = database.SetTransaction(ctx, transaction)

// 	user := &model.User{Username: "reflect15@test.com", Firstname: "Test", Lastname: "Reflection", DOB: time.Now()}
// 	err := Insert(ctx, user)

// 	assert.True(t, user.ID != 0, "ID not set after Insert")
// 	assert.NoError(t, err, "Error during Insert")

// 	transaction.Rollback()
// }

// func TestDelete(t *testing.T) {

// 	ctx := context.Background()
// 	transaction := database.NewTransaction(ctx)
// 	ctx = database.SetTransaction(ctx, transaction)

// 	user := &model.User{ID: 169}
// 	err := Delete(ctx, user)
// 	assert.NoError(t, err, "Error during Delete")

// 	transaction.Rollback()
// }

// func TestCount(t *testing.T) {

// 	ctx := context.Background()
// 	transaction := database.NewTransaction(ctx)
// 	ctx = database.SetTransaction(ctx, transaction)

// 	user := &model.User{Firstname: "Joe"}
// 	count, err := Count(ctx, user)
// 	assert.NoError(t, err, "Error during Count")
// 	assert.True(t, count > 0, "Count returned 0")

// 	transaction.Rollback()
// }

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

func TestGetEntityData(t *testing.T) {
	expectedEntityData := &entityData{}
	expectedEntityData.Name = "public.testentity"
	expectedEntityData.Attributes = map[string]interface{}{
		"id":     1,
		"name":   "Joe",
		"column": "Age",
	}
	entityData := getEntityData(testEntity)
	assert.Equal(t, expectedEntityData, entityData)
}

func TestGetEntityDataEmptyData(t *testing.T) {
	entity := &TestEntity{}
	expectedEntityData := &entityData{}
	expectedEntityData.Name = "public.testentity"
	expectedEntityData.Attributes = map[string]interface{}{}
	entityData := getEntityData(entity)
	assert.Equal(t, expectedEntityData, entityData)
}
