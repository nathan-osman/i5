package db

import (
	"errors"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

var (
	usersBucket = []byte("Users")

	errInvalidUsernamePassword = errors.New("invalid username or password")
)

// CreateUser creates a new user in the database.
func (d *DB) CreateUser(username, password string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket(usersBucket)
		if err != nil {
			return err
		}
		return b.Put([]byte(username), h)
	})
}

// LoginUser verifies the specified login credentials.
func (d *DB) LoginUser(username, password string) error {
	var hashedPassword []byte
	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(usersBucket)
		if b == nil {
			return nil
		}
		hashedPassword = b.Get([]byte(username))
		return nil
	})
	if hashedPassword == nil {
		return errInvalidUsernamePassword
	}
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return errInvalidUsernamePassword
	}
	return nil
}
