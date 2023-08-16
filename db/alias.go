package db

import "gorm.io/gorm"

type Alias struct {
	gorm.Model
	Name    string `gorm:"not null;unique"`
	Content string `gorm:"not null"`
}

func CreateAlias(name string, content string) {
	db.Create(&Alias{
		Name:    name,
		Content: content,
	})
}

func DeleteAlias(name string) {
	db.Unscoped().Delete(&Alias{}, "name = ?", name)
}

func GetAlias(name string) *Alias {
	var alias Alias
	db.First(&alias, "name = ?", name)
	if alias.Name == "" {
		return nil
	}
	return &alias
}
