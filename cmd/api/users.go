package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllUsers Return all users
// @Summary       Returns all users
// @Description   Return all users
// @Tags          Users
// @Produce       json
// @Success       200             {object} []database.User
// @Router        /api/v1/users   [get]
func (app *app) getAllUsers(c *gin.Context) {
	users, err := app.models.Users.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}
