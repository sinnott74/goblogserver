package model

import (
	"context"
	"testing"
	"time"

	"github.com/sinnott74/goblogserver/orm"
	"github.com/stretchr/testify/assert"
)

func BenchmarkSaveUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		transaction, ctx := orm.NewTransaction(context.Background())

		user := &User{Username: "reflect16@test.com", Firstname: "Test", Lastname: "Reflection", DOB: time.Now().UTC().Round(time.Microsecond)}
		err := orm.Save(ctx, user)
		assert.NoError(b, err, "Error during initial save")
		assert.True(b, user.ID != 0, "ID not set after Insert")

		user.Lastname = "Save"
		err = orm.Save(ctx, user)
		assert.NoError(b, err, "Error during modify save")

		userByID := &User{ID: user.ID}
		err = orm.Get(ctx, userByID)
		assert.NoError(b, err, "Error reading a User back")
		assert.Equal(b, user, userByID, "User saved is not the same as user read")
		transaction.Rollback()
	}
}

func TestSaveUser(t *testing.T) {
	transaction, ctx := orm.NewTransaction(context.Background())
	defer transaction.Rollback()

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
}

func TestInsertUser(t *testing.T) {

	transaction, ctx := orm.NewTransaction(context.Background())
	defer transaction.Rollback()

	user := &User{Username: "reflect16@test.com", Firstname: "Test", Lastname: "Reflection", DOB: time.Now().UTC().Round(time.Microsecond)}
	err := orm.Insert(ctx, user)

	assert.NoError(t, err, "Error during Insert")
	assert.True(t, user.ID != 0, "ID not set after Insert")

	userByID := &User{ID: user.ID}
	err = orm.Get(ctx, userByID)
	assert.NoError(t, err, "Error reading a User back")
	assert.Equal(t, user, userByID, "User inserted is not the same as user read")
}

func TestInsertUserWithUsernameTaken(t *testing.T) {

	transaction, ctx := orm.NewTransaction(context.Background())
	defer transaction.Rollback()

	user := &User{Username: "reflect15@test.com", Firstname: "Test", Lastname: "Reflection", DOB: time.Now()}
	err := orm.Insert(ctx, user)

	assert.Error(t, err, "Username taken error not returned")
}

func TestUpdateUser(t *testing.T) {

	transaction, ctx := orm.NewTransaction(context.Background())
	defer transaction.Rollback()

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
}

func TestDeleteUser(t *testing.T) {

	transaction, ctx := orm.NewTransaction(context.Background())
	defer transaction.Rollback()

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
}
