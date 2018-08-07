package auth

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sinnott74/goblogserver/env"
	"github.com/sinnott74/goblogserver/model"
	"github.com/sinnott74/goblogserver/orm"
)

// ErrorAuth is an Authentication Error
var ErrorAuth = errors.New("Authentication Errror")

var authHeader = "Authorization"

// UserToken authentication token given to a client containing the JWT token & the user details
type UserToken struct {
	model.User
	Token string `json:"token"`
}

type jwtClaims struct {
	model.User
	jwt.StandardClaims
}

// Global JWT Secret
var secret = env.JWTSecret()

// Login logs a user in & generates a JWT token
func Login(ctx context.Context, username string, password string) (*UserToken, error) {
	authenticated := model.Authenticate(ctx, username, password)
	if !authenticated {
		return nil, ErrorAuth
	}
	user := &model.User{Username: username}
	err := orm.SelectOne(ctx, user)
	if err != nil {
		return nil, err
	}
	jwtClaims := jwtClaims{*user, jwt.StandardClaims{ExpiresAt: time.Now().AddDate(0, 0, 7).Unix()}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}
	return &UserToken{User: *user, Token: tokenString}, nil
}

//Authenticate verifies that both the tokenString is valid & the user it represents
func Authenticate(ctx context.Context, tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		jwtClaims, ok := token.Claims.(*jwtClaims)
		if !ok {
			return false, ErrorAuth
		}
		user := &model.User{Username: jwtClaims.Username}
		err := orm.SelectOne(ctx, user)
		if err != nil {
			return false, ErrorAuth
		}
		if !reflect.DeepEqual(&jwtClaims.User, user) {
			return false, ErrorAuth
		}
		return secret, nil
	})

	return token.Valid, err
}

// Middleware which authenticates using a JWT header
func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get(authHeader)
		authenticated, err := Authenticate(ctx, token)
		if !authenticated || err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
