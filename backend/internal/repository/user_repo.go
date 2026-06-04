package repository

import (
	"cowork/internal/db"
	"cowork/internal/model"

	"gorm.io/gorm"
)

// CreateUser inserts a new user into the database.
func CreateUser(user *model.User) error {
	return db.DB.Create(user).Error
}

// FindByUsername finds a user by username. Returns nil if not found.
func FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email. Returns nil if not found.
func FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := db.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by primary key ID.
func FindByID(id uint) (*model.User, error) {
	var user model.User
	err := db.DB.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user record.
func UpdateUser(user *model.User) error {
	return db.DB.Save(user).Error
}
