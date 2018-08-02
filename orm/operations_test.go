package orm

import (
	"testing"

	"github.com/iancoleman/strcase"

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

func TestGetEntityData(t *testing.T) {
	config := Config{
		ToDBMapperFunc: strcase.ToSnake,
	}
	expectedEntityData := &entityData{}
	expectedEntityData.Name = "public.testentity"
	expectedEntityData.Attributes = map[string]interface{}{
		"id":     1,
		"name":   "Joe",
		"column": "Age",
	}
	entityData := getEntityData(testEntity, config)
	assert.Equal(t, expectedEntityData, entityData)
}

func TestGetEntityDataEmptyData(t *testing.T) {
	config := Config{
		ToDBMapperFunc: strcase.ToSnake,
	}
	entity := &TestEntity{}
	expectedEntityData := &entityData{}
	expectedEntityData.Name = "public.testentity"
	expectedEntityData.Attributes = map[string]interface{}{}
	entityData := getEntityData(entity, config)
	assert.Equal(t, expectedEntityData, entityData)
}
