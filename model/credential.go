package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sinnott74/goblogserver/orm"

	"golang.org/x/crypto/bcrypt"
)

// Credential entity
type Credential struct {
	ID        int64     `json:"id"`
	Password  string    `json:"password"`
	Active    bool      `json:"active"`
	CreatedOn time.Time `json:"created_on"`
	UserID    int64     `json:"user_id"`
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

// deactivePreviousCredential deactives all credentials for this user
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

// comparePassword checks the given password against the current credentials encryped password
func (c *Credential) comparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))
	return err == nil
}

// Authenticate validates that the given username & password pairing represents an user
func Authenticate(ctx context.Context, username string, password string) bool {
	credential, err := readActiveCredentialByUserName(ctx, username)
	if err != nil {
		panic(err)
	}
	return credential.comparePassword(password)
}

// readActiveCredentialByUserName reads the active credential for the given username
func readActiveCredentialByUserName(ctx context.Context, username string) (*Credential, error) {
	credential := &Credential{}
	rows, err := orm.Exec(ctx, "SELECT * FROM credential WHERE active=true AND user_id = (SELECT id FROM public.user WHERE username = $1)", username)
	if err != nil {
		return nil, err
	}
	err = rows.StructScan(credential)
	if err != nil {
		return nil, err
	}
	return credential, nil
}
