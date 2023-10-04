package db

import (
	"time"

	"gorm.io/gorm"
)

type RepeatedCommand struct {
	gorm.Model
	JID      string    `gorm:"not null;column:jid"`
	User     string    `gorm:"not null"`
	Content  string    `gorm:"not null"`
	Repeat   string    `gorm:"not null"`
	NextDate time.Time `gorm:"not null"`
}

func CreateRepeatedCommand(jid string, user string, content string, repeat string, nextDate time.Time) (*RepeatedCommand, error) {
	data := &RepeatedCommand{
		JID:      jid,
		User:     user,
		Content:  content,
		Repeat:   repeat,
		NextDate: nextDate,
	}

	err := db.Create(data).Error
	return data, err
}

func ClearRepeatedCommands(jid string, user string) (int64, error) {
	result := db.Unscoped().Where("jid = ? AND user = ?", jid, user).Delete(&RepeatedCommand{})
	return result.RowsAffected, result.Error
}

func UpdateNextDate(id uint, nextDate time.Time) error {
	return db.Model(&RepeatedCommand{}).Where("id = ?", id).Update("next_date", nextDate).Error
}

func GetRepeatedCommandsToday() []RepeatedCommand {
	var messages []RepeatedCommand

	now := time.Now()
	db.Where("next_date <= ?", now).Find(&messages)

	return messages
}
