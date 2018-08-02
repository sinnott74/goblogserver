package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
)

// const MultipleRecordsFound = errors.New("Multiple Records Found")
// const SingleRecordExpected = errors.New("Single Records Expected")

type EntityID struct {
	ID int64
}

// ExecuteQueryRow is the central function for execting queries
func ExecuteQueryRow(ctx context.Context, query string, values ...interface{}) *sqlx.Row {
	t := getTransaction(ctx)
	if config.Debug {
		fmt.Printf("%s - %s - %+v\n", t.ID(), query, values)
	}
	return t.Tx().QueryRowxContext(ctx, query, values...)
}

// Save saves an entity
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
	entityData := getEntityData(entity, config)
	keys := getOrderedKeys(entityData.Attributes)
	query := insertQuery(entityData.Name, keys)
	var values []interface{}
	for i := 0; i < len(keys); i++ {
		values = append(values, entityData.Attributes[keys[i]])
	}
	var lastID int64
	err = ExecuteQueryRow(ctx, query, values...).Scan(&lastID)
	if err != nil {
		return err
	}
	setID(entity, lastID)

	err = callHook(entity, "PostInsert", ctx)
	return err
}

// Get reads a database row by ID
func Get(ctx context.Context, entity interface{}) error {
	entityData := getEntityData(entity, config)
	query := getQuery(entityData.Name)
	ID := getID(entity)
	return ExecuteQueryRow(ctx, query, ID).StructScan(entity)
}

// Update updates the entites specified by the where with the value spcified in the set
func Update(ctx context.Context, set interface{}, where interface{}) error {
	setEntityData := getEntityData(set, config)
	whereEntityData := getEntityData(where, config)
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
	err := ExecuteQueryRow(ctx, query, values...).Scan()
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
	entityData := getEntityData(entity, config)
	query := deleteByIDQuery(entityData.Name)
	ID := getID(entity)
	err = ExecuteQueryRow(ctx, query, ID).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return callHook(entity, "PostDelete", ctx)
}

// SelectAll performs an SQL select query
func SelectAll(ctx context.Context, entities interface{}, where interface{}) error {
	whereEntityData := getEntityData(where, config)
	whereKeys := getOrderedKeys(whereEntityData.Attributes)
	query := selectQuery(whereEntityData.Name, nil, whereKeys)
	var values []interface{}
	for _, whereKey := range whereKeys {
		values = append(values, whereEntityData.Attributes[whereKey])
	}
	t := getTransaction(ctx)
	if config.Debug {
		fmt.Printf("%s - %s - %+v\n", t.ID(), query, values)
	}
	err := t.Tx().Select(entities, query, values...)
	return err
}

// SelectOne performs an SQL select query
func SelectOne(ctx context.Context, entity interface{}) error {
	entityData := getEntityData(entity, config)
	entityKeys := getOrderedKeys(entityData.Attributes)
	query := selectQuery(entityData.Name, nil, entityKeys)
	var values []interface{}
	for _, entityKey := range entityKeys {
		values = append(values, entityData.Attributes[entityKey])
	}
	return ExecuteQueryRow(ctx, query, values...).StructScan(entity)
}

// Count Counts the number of rows with the given where values
func Count(ctx context.Context, where interface{}) (int64, error) {
	whereEntityData := getEntityData(where, config)
	whereKeys := getOrderedKeys(whereEntityData.Attributes)
	query := countQuery(whereEntityData.Name, whereKeys)
	var values []interface{}
	for i := 0; i < len(whereKeys); i++ {
		values = append(values, whereEntityData.Attributes[whereKeys[i]])
	}
	var count int64
	error := ExecuteQueryRow(ctx, query, values...).Scan(&count)
	if error != nil {
		return 0, error
	}
	return count, nil
}

type entityData struct {
	Name       string
	Attributes map[string]interface{}
}

// getEntityData uses reflection to pull out the entity's data. i.e its name, columns names & values
func getEntityData(entity interface{}, config Config) *entityData {
	entityData := &entityData{}
	typ := reflect.TypeOf(entity).Elem()
	val := reflect.ValueOf(entity).Elem()
	entityData.Name = "public." + strings.ToLower(typ.Name())
	entityData.Attributes = make(map[string]interface{})
	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		columnName := strcase.ToSnake(fieldType.Name)
		// columnName := getColumnName(fieldType, config)
		fieldVal := val.Field(i)
		if fieldVal.Interface() != reflect.Zero(fieldVal.Type()).Interface() {
			entityData.Attributes[columnName] = fieldVal.Interface()
		}
	}
	return entityData
}

// getColumnName gets the column name associated with the given struct field
func getColumnName(field reflect.StructField, config Config) string {
	if config.ToDBMapperFunc != nil {
		return strcase.ToSnake(field.Name)
	}
	return field.Name
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
