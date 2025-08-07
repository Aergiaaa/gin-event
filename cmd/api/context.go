package main

import (
	"github.com/Aergiaaa/gin-event/internal/database"

	"github.com/gin-gonic/gin"
)

func (app *app) getUserFromContext(c *gin.Context) *database.User {
	var emptyUser database.User

	ctxUser, exist := c.Get("user")
	if !exist {
		return &emptyUser
	}

	user, ok := ctxUser.(*database.User)
	if !ok {
		return &emptyUser
	}

	return user
}

