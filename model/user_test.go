package model

import (
	"context"
	"testing"
	"time"

	"github.com/sinnott74/goblogserver/database"
	"github.com/sinnott74/goblogserver/orm"
	"github.com/stretchr/testify/assert"
)

func TestSaveUser(t *testing.T) {
	ctx := context.Background()
	transaction := database.NewTransaction(ctx)
	ctx = database.SetTransaction(ctx, transaction)

	user := &User{Username: "reflect16@test.com", Firstname: "Test", Lastname: "Reflection", DOB: time.Now().UTC().Round(time.Microsecond)}
	err := orm.Save(ctx, user)
	assert.NoError(t, err, "Error during initial save")
	assert.True(t, user.ID != 0, "ID not set after Insert")

	user.Lastname = "Save"
	err = orm.Save(ctx, user)
	assert.NoError(t, err, "Error during modify save")

	userByID := &User{ID: user.ID}
	err = orm.Get(ctx, userByID)
	assert.NoError(t, err, "Error reading a User back")
	assert.Equal(t, user, userByID, "User saved is not the same as user read")
	transaction.Rollback()
}

func TestInsertUser(t *testing.T) {

	ctx := context.Background()
	transaction := database.NewTransaction(ctx)
	ctx = database.SetTransaction(ctx, transaction)

	user := &User{Username: "reflect16@test.com", Firstname: "Test", Lastname: "Reflection", DOB: time.Now().UTC().Round(time.Microsecond)}
	err := orm.Insert(ctx, user)

	assert.NoError(t, err, "Error during Insert")
	assert.True(t, user.ID != 0, "ID not set after Insert")

	userByID := &User{ID: user.ID}
	err = orm.Get(ctx, userByID)
	assert.NoError(t, err, "Error reading a User back")
	assert.Equal(t, user, userByID, "User inserted is not the same as user read")
	transaction.Rollback()
}

func TestInsertUserWithUsernameTaken(t *testing.T) {

	ctx := context.Background()
	transaction := database.NewTransaction(ctx)
	ctx = database.SetTransaction(ctx, transaction)

	user := &User{Username: "reflect15@test.com", Firstname: "Test", Lastname: "Reflection", DOB: time.Now()}
	err := orm.Insert(ctx, user)

	assert.Error(t, err, "Username taken error not returned")
	transaction.Rollback()
}

func TestUpdateUser(t *testing.T) {

	ctx := context.Background()
	transaction := database.NewTransaction(ctx)
	ctx = database.SetTransaction(ctx, transaction)

	user := &User{Username: "reflect16@test.com", Firstname: "Test", Lastname: "Reflection", DOB: time.Now().UTC().Round(time.Microsecond)}
	err := orm.Insert(ctx, user)
	assert.NoError(t, err, "Error during Insert")

	userByID := &User{ID: user.ID}
	userLastname := &User{Lastname: "Modified"}

	err = orm.Update(ctx, userLastname, userByID)
	assert.NoError(t, err, "Error during modify of lastname")

	expectedUser := &User{ID: user.ID, Username: user.Username, Firstname: user.Firstname, Lastname: userLastname.Lastname, DOB: user.DOB}
	err = orm.Get(ctx, userByID)
	assert.NoError(t, err, "Error reading a User back")
	assert.Equal(t, expectedUser, userByID, "User modifed is not the same as user read")
	transaction.Rollback()
}

func TestDeleteUser(t *testing.T) {

	ctx := context.Background()
	transaction := database.NewTransaction(ctx)
	ctx = database.SetTransaction(ctx, transaction)

	user := &User{Username: "reflect16@test.com", Firstname: "Test", Lastname: "Reflection", DOB: time.Now().UTC().Round(time.Microsecond)}
	err := orm.Insert(ctx, user)
	assert.NoError(t, err, "Error during Insert")

	expectedCountBeforeDelete := int64(1)
	userByID := &User{ID: user.ID}
	count, err := orm.Count(ctx, userByID)
	assert.NoError(t, err, "Error during count before delete")
	assert.Equal(t, expectedCountBeforeDelete, count, "User count incorrect")

	err = orm.Delete(ctx, userByID)
	assert.NoError(t, err, "Error during Delete")

	expectedCountAfterDelete := int64(0)
	count, err = orm.Count(ctx, userByID)
	assert.NoError(t, err, "Error during count after delete")
	assert.Equal(t, expectedCountAfterDelete, count, "User count incorrect")
	transaction.Rollback()
}
