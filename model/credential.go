package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sinnott74/goblogserver/orm"

	"github.com/sinnott74/goblogserver/database"

	"golang.org/x/crypto/bcrypt"
)

// Credential entity
type Credential struct {
	ID        int64
	Password  string
	Active    bool
	CreatedOn time.Time
	UserID    int64
}

// PreInsert hook to encrypt password & deactivate previous passwords
func (c *Credential) PreInsert(ctx context.Context) error {
	err := c.deactivePreviousCredential(ctx)
	if err != err {
		return err
	}
	c.CreatedOn = time.Now().UTC()
	c.Active = true
	err = c.encrypt()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *Credential) deactivePreviousCredential(ctx context.Context) error {
	prevActiveCred := &Credential{UserID: c.UserID, Active: true}
	set := &Credential{Active: true}
	return orm.Update(ctx, set, prevActiveCred)
}

// encrypts the given string
func (c *Credential) encrypt() error {
	if len(c.Password) == 0 {
		return errors.New("Password should not be empty")
	}
	bytePassword := []byte(c.Password)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.Password = string(passwordHash)
	return nil
}

func (c *Credential) comparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))
	return err == nil
}

// Authenticate validates that the given username & password pairing represents an user
func (c *Credential) Authenticate(ctx context.Context, username string, password string) bool {
	credential, err := readActiveCredentialByUserName(ctx, username)
	if err != nil {
		panic(err)
	}
	return credential.comparePassword(password)
}

func readActiveCredentialByUserName(ctx context.Context, username string) (Credential, error) {
	var credential Credential
	t := ctx.Value(database.TransactionKey).(database.Transaction)
	err := t.Tx().QueryRowContext(ctx, "SELECT * FROM credential WHERE user_id = (SELECT id FROM public.user WHERE username = ?)", username).Scan(&credential)
	if err != nil {
		return credential, err
	}
	return credential, nil
}
