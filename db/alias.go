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

	gorm := db.Create(data)
	return data, gorm.Error
}

func UpsertAlias(name string, content string) error {
	alias, err := GetAlias(name)
	if err != nil {
		return err
	}

	if alias == nil {
		_, err := CreateAlias(name, content)
		return err
	}

	gorm := db.Model(&alias).Update("content", content)

	return gorm.Error
}

func DeleteAlias(name string) error {
	gorm := db.Unscoped().Delete(&Alias{}, "name = ?", name)
	return gorm.Error
}

func GetAlias(name string) (*Alias, error) {
	var alias Alias

	gorm := db.First(&alias, "name = ?", name)
	if alias.Name == "" {
		return nil, gorm.Error
	}

	return &alias, gorm.Error
}
