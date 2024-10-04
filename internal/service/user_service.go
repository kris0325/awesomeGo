package service

import (
	"awesomeGo/internal/model"
	"awesomeGo/internal/repository"
)

// CreateUser creates a new user
func CreateUser(user *model.User) error {
	return repository.CreateUser(user)
}

// GetUsers retrieves all users
func GetUsers() ([]model.User, error) {
	return repository.FindAllUsers()
}

// GetUserByID retrieves a user by ID
func GetUserByID(id string) (model.User, error) {
	return repository.FindUserByID(id)
}

// UpdateUser updates a user
func UpdateUser(user *model.User) error {
	return repository.UpdateUser(user)
}

// DeleteUser deletes a user by ID
func DeleteUser(id string) error {
	return repository.DeleteUser(id)
}
