package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sinnott74/goblogserver/orm"
)

// User entity
type User struct {
	ID        int64
	Username  string
	Firstname string
	Lastname  string
	DOB       time.Time
}

// PreInsert is called before a user is inserted into the User table
func (u *User) PreInsert(ctx context.Context) error {
	fmt.Println(u.Firstname + u.Lastname)
	if !u.isUserNameAvailable(ctx) {
		return errors.New("Username taken")
	}
	return nil
}

// isUserNameAvailable checks is this user's username is available
func (u *User) isUserNameAvailable(ctx context.Context) bool {
	userWithSameName := &User{Username: u.Username}
	count, err := orm.Count(ctx, userWithSameName)
	return count == 0 && err == nil
}