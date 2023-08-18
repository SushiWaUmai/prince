package db

import "gorm.io/gorm"

type Alias struct {
	gorm.Model
	Name    string `gorm:"not null;unique"`
	Content string `gorm:"not null"`
}

func CreateAlias(name string, content string) (*Alias, error) {
	data := &Alias{
		Name:    name,
		Content: content,
	}

	err := db.Create(data).Error
	return data, err
}

func UpsertAlias(name string, content string) error {
	alias := GetAlias(name)

	if alias == nil {
		_, err := CreateAlias(name, content)
		return err
	}

	err := db.Model(&alias).Update("content", content).Error

	return err
}

func DeleteAlias(name string) error {
	return db.Unscoped().Delete(&Alias{}, "name = ?", name).Error
}

func GetAlias(name string) *Alias {
	var alias Alias

	db.First(&alias, "name = ?", name)
	if alias.Name == "" {
		return nil
	}

	return &alias
}
