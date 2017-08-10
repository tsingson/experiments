package models

// User model
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// A slice of users model collection
type Users []User

// An instance of users struct
var users Users

// Get all users
func (self User) All() Users {
	return users
}

// Create a new user
func (self User) Create(name string, age int) string {
	users = append(users, User{name, age})

	return "User has been created"
}

// Clear users
func (self User) Clear() {
	users = Users{}
}
