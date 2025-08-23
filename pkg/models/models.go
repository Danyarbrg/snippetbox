package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// add a new ErrInvalidCredentials error to use this if a user
	// tries to login with an incorrect email address or password
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// add a new ErrDuplicateEmail error. We will use this if a user trues to signup
	// with an email addres thats already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type Snippet struct {
	ID		int
	Title	string
	Content	string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID				int
	Name			string
	Email			string
	HashedPassword	[]byte
	Created			time.Time
	Active			bool
}