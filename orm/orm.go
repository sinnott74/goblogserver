package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/sinnott74/goblogserver/database"
)

// const MultipleRecordsFound = errors.New("Multiple Records Found")
// const SingleRecordExpected = errors.New("Single Records Expected")

type EntityID struct {
	ID int64
}

// executeQueryRow is the central function for execting queries
func executeQueryRow(ctx context.Context, query string, values ...interface{}) *sqlx.Row {
	t := database.GetTransaction(ctx)
	fmt.Printf("%s - %s - %+v\n", t.ID(), query, values)
	return t.Tx().QueryRowxContext(ctx, query, values...)
}

// Saves an entity
func Save(ctx context.Context, entity interface{}) error {
	err := callHook(entity, "PreSave", ctx)
	if err != nil {
		return err
	}
	id := getID(entity)
	if id == 0 {
		err = Insert(ctx, entity)
	} else {
		entityID := &EntityID{id}
		err = Update(ctx, entity, entityID)
	}
	if err != nil {
		return err
	}
	return callHook(entity, "PostSave", ctx)
}

// Insert inserts an entity into the database & sets the returned ID onto the entity
func Insert(ctx context.Context, entity interface{}) error {
	err := callHook(entity, "PreInsert", ctx)
	if err != nil {
		return err
	}
	entityData := getEntityData(entity)
	keys := getOrderedKeys(entityData.Attributes)
	query := insertQuery(entityData.Name, keys)
	var values []interface{}
	for i := 0; i < len(keys); i++ {
		values = append(values, entityData.Attributes[keys[i]])
	}
	var lastID int64
	err = executeQueryRow(ctx, query, values...).Scan(&lastID)
	if err != nil {
		return err
	}
	setID(entity, lastID)

	err = callHook(entity, "PostInsert", ctx)
	return err
}

// Get reads a database row by ID
func Get(ctx context.Context, entity interface{}) error {
	entityData := getEntityData(entity)
	query := getQuery(entityData.Name)
	ID := getID(entity)
	return executeQueryRow(ctx, query, ID).StructScan(entity)
}

// Update updates the entites specified by the where with the value spcified in the set
func Update(ctx context.Context, set interface{}, where interface{}) error {
	setEntityData := getEntityData(set)
	whereEntityData := getEntityData(where)
	setKeys := getOrderedKeys(setEntityData.Attributes)
	whereKeys := getOrderedKeys(whereEntityData.Attributes)
	var values []interface{}
	for _, setKey := range setKeys {
		values = append(values, setEntityData.Attributes[setKey])
	}
	for _, whereKey := range whereKeys {
		values = append(values, whereEntityData.Attributes[whereKey])
	}
	query := updateQuery(setEntityData.Name, setKeys, whereKeys)
	err := executeQueryRow(ctx, query, values...).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

// Delete removes an row from the database. Rows are deleted by ID
func Delete(ctx context.Context, entity interface{}) error {
	err := callHook(entity, "PreDelete", ctx)
	if err != nil {
		return err
	}
	entityData := getEntityData(entity)
	query := deleteByIDQuery(entityData.Name)
	ID := getID(entity)
	err = executeQueryRow(ctx, query, ID).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return callHook(entity, "PostDelete", ctx)
}

// SelectAll performs an SQL select query
func SelectAll(ctx context.Context, entities interface{}, where interface{}) error {
	whereEntityData := getEntityData(where)
	whereKeys := getOrderedKeys(whereEntityData.Attributes)
	query := selectQuery(whereEntityData.Name, nil, whereKeys)
	var values []interface{}
	for _, whereKey := range whereKeys {
		values = append(values, whereEntityData.Attributes[whereKey])
	}
	t := database.GetTransaction(ctx)
	fmt.Printf("%s - %s - %+v\n", t.ID(), query, values)
	err := t.Tx().Select(entities, query, values...)
	return err
}

// Count Counts the number of rows with the given where values
func Count(ctx context.Context, where interface{}) (int64, error) {
	whereEntityData := getEntityData(where)
	whereKeys := getOrderedKeys(whereEntityData.Attributes)
	query := countQuery(whereEntityData.Name, whereKeys)
	var values []interface{}
	for i := 0; i < len(whereKeys); i++ {
		values = append(values, whereEntityData.Attributes[whereKeys[i]])
	}
	var count int64
	error := executeQueryRow(ctx, query, values...).Scan(&count)
	if error != nil {
		return 0, error
	}
	return count, nil
}

// createGetQuery creates a select by id query
// e.g. SELECT * FROM tablename;
func getQuery(tableName string) string {
	return selectQuery(tableName, nil, []string{"id"})
}

// deleteByIDQuery creates a Delete query which delete by a rows ID field
// e.g. DELETE FROM test WHERE id=$1
func deleteByIDQuery(tableName string) string {
	return deleteQuery(tableName, []string{"id"})
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

type entityData struct {
	Name       string
	Attributes map[string]interface{}
}

// getEntityData uses reflection to pull out the entity's data. i.e its name, columns names & values
func getEntityData(entity interface{}) *entityData {
	entityData := &entityData{}
	typ := reflect.TypeOf(entity).Elem()
	val := reflect.ValueOf(entity).Elem()
	entityData.Name = "public." + strings.ToLower(typ.Name())
	entityData.Attributes = make(map[string]interface{})
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		columnName := strcase.ToSnake(fieldType.Name)
		fieldVal := val.Field(i)
		if fieldVal.Interface() != reflect.Zero(fieldVal.Type()).Interface() {
			entityData.Attributes[columnName] = fieldVal.Interface()
		}
	}
	return entityData
}

// getID Reads an entities ID field
func getID(entity interface{}) int64 {
	return reflect.ValueOf(entity).Elem().FieldByName("ID").Int()
}

// setID sets the give ID as the ID of the struct
func setID(entity interface{}, ID int64) {
	reflect.ValueOf(entity).Elem().FieldByName("ID").SetInt(ID)
}

// getOrderedKeys returns a sorted list of the map's key
// removeID specifies whether the key id should be removed from the final list
func getOrderedKeys(m map[string]interface{}) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// callHook calls the given hook on the given entity
func callHook(entity interface{}, hookname string, args ...interface{}) error {
	hook := reflect.ValueOf(entity).MethodByName(hookname)

	if !hook.IsValid() {
		return nil
	}

	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	res := hook.Call(inputs)

	if err, ok := res[0].Interface().(error); ok && err != nil {
		return err
	}
	return nil
}

// *** TODO ***
// FindAll
// FindAtMostOne
// FindOne
// Before/After Hooks
// Insert multiple
