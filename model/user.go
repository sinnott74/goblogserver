package model

import (
	"context"
	"errors"
	"time"

	"github.com/sinnott74/goblogserver/orm"
)

// User entity
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	DOB       time.Time `json:"dob"`
}

// PreInsert is called before a user is inserted into the User table
func (u *User) PreInsert(ctx context.Context) error {
	usernameAvailable, err := u.isUserNameAvailable(ctx)
	if err != nil {
		return err
	}
	if !usernameAvailable {
		return errors.New("Username taken")
	}
	return nil
}

// isUserNameAvailable checks is this user's username is available
func (u *User) isUserNameAvailable(ctx context.Context) (bool, error) {
	userWithSameName := &User{Username: u.Username}
	count, err := orm.Count(ctx, userWithSameName)
	return count == 0, err
}
