package repository

import (
	"awesomeGo/internal/config"
	"awesomeGo/internal/model"
)

// CreateUser creates a new user in the database
func CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

// FindAllUsers retrieves all users from the database
func FindAllUsers() ([]model.User, error) {
	var users []model.User
	err := config.DB.Find(&users).Error
	return users, err
}

// FindUserByID retrieves a user by ID
func FindUserByID(id string) (model.User, error) {
	var user model.User
	err := config.DB.First(&user, id).Error
	return user, err
}

// UpdateUser updates an existing user
func UpdateUser(user *model.User) error {
	return config.DB.Save(user).Error
}

// DeleteUser deletes a user by ID
func DeleteUser(id string) error {
	return config.DB.Delete(&model.User{}, id).Error
}
