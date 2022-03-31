package status

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionUsername = "username"

	errInvalidUsernamePassword = "invalid username or password"
	errUnauthorized            = "you do not have permission to access that page"
)

type authLoginParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *Status) authLogin(c *gin.Context) {
	params := &authLoginParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		failure(c, http.StatusBadRequest, err.Error())
		return
	}
	var password []byte
	s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(usersBucket)
		if b == nil {
			return nil
		}
		password = b.Get([]byte(params.Username))
		return nil
	})
	if password == nil {
		failure(c, http.StatusUnauthorized, errInvalidUsernamePassword)
		return
	}
	if err := bcrypt.CompareHashAndPassword(password, []byte(params.Password)); err != nil {
		failure(c, http.StatusUnauthorized, errInvalidUsernamePassword)
		return
	}
	session := sessions.Default(c)
	session.Set(sessionUsername, params.Username)
	if err := session.Save(); err != nil {
		failure(c, http.StatusInternalServerError, err.Error())
		return
	}
	success(c)
}

func requireLogin(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get(sessionUsername) == nil {
		failure(c, http.StatusUnauthorized, errUnauthorized)
		c.Abort()
		return
	}
}
