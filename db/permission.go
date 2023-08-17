package db

import "gorm.io/gorm"

type UserPermission struct {
	gorm.Model
	UserID     string `gorm:"not null;unique;column:user_id"`
	Permission string `gorm:"not null;column:permission"`
}

func GetUserPermission(userId string) (*UserPermission, error) {
	// Gets User Permission, if it doesn't exist create one with Permission NONE
	var userPerm UserPermission
	err := db.FirstOrInit(&userPerm, UserPermission{UserID: userId}).Error
	if err != nil {
		return nil, err
	}

	if userPerm.Permission == "" {
		userPerm.Permission = "NONE"
	}
	err = db.Save(&userPerm).Error

	return &userPerm, nil
}

func UpdateUserPermission(userId string, permission string) error {
	return db.Model(&UserPermission{}).Where("user_id = ?", userId).Update("permission", permission).Error
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
