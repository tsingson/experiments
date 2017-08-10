package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tsingson/test/goconvey/models"
)

// Users controller struct
type UsersController struct {
	User models.User
}

// Create a new user
func (self *UsersController) Create(c *gin.Context) {
	name := c.PostForm("name")
	age, _ := strconv.Atoi(c.PostForm("age"))

	c.JSON(201, gin.H{
		"message": self.User.Create(name, age),
	})
}

// Get all users
func (self *UsersController) Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": self.User.All(),
	})
}
