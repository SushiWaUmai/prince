package db

import "gorm.io/gorm"

type UserPermission struct {
	gorm.Model
	UserID     string `gorm:"not null;unique"`
	Permission string `gorm:"not null"`
}

func GetUserPermission(userId string) UserPermission {
	// Gets User Permission, if it doesn't exist create one with Permission NONE
	var userPerm UserPermission
	db.FirstOrInit(&userPerm, UserPermission{UserID: userId})
	if userPerm.Permission == "" {
		userPerm.Permission = "NONE"
	}
	db.Save(&userPerm)
	return userPerm
}

func UpsertPermission(userId string, permission string) UserPermission {
	userPerm := UserPermission{
		UserID:     userId,
		Permission: permission,
	}
	db.Save(&userPerm)
	return userPerm
}

func ComparePermission(perm string, command string) bool {
	// NONE, USER, ADMIN, OP
	switch perm {
	case "NONE":
		return false
	case "USER":
		if command == "ADMIN" || command == "OP" {
			return false
		}
	case "ADMIN":
		if command == "OP" {
			return false
		}
	}
	return true
}
