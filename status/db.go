package status

import (
	"path"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

var (
	usersBucket = []byte("Users")
)

func openDB(storageDir string) (*bolt.DB, error) {
	return bolt.Open(
		path.Join(storageDir, "i5.db"),
		0600,
		nil,
	)
}

// CreateUser initializes the database and creates a new user.
func CreateUser(storageDir, username, password string) error {
	d, err := openDB(storageDir)
	if err != nil {
		return err
	}
	defer d.Close()
	h, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}
	return d.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket(usersBucket)
		if err != nil {
			return err
		}
		return b.Put([]byte(username), h)
	})
}
