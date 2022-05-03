package status

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	sessionUsername = "username"

	errUnauthorized = "you do not have permission to access that page"
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
	if err := s.db.LoginUser(params.Username, params.Password); err != nil {
		failure(c, http.StatusUnauthorized, err.Error())
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

func (s *Status) authLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(sessionUsername)
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
