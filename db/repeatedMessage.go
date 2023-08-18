package db

import (
	"time"

	"gorm.io/gorm"
)

type RepeatedMessage struct {
	gorm.Model
	JID      string    `gorm:"not null;column:jid"`
	User     string    `gorm:"not null"`
	Message  string    `gorm:"not null"`
	Repeat   string    `gorm:"not null"`
	NextDate time.Time `gorm:"not null"`
}

func CreateRepeatedMessage(jid string, user string, message string, repeat string, nextDate time.Time) (*RepeatedMessage, error) {
	data := &RepeatedMessage{
		JID:      jid,
		User:     user,
		Message:  message,
		Repeat:   repeat,
		NextDate: nextDate,
	}

	err := db.Create(data).Error
	return data, err
}

func ClearRepeatedMessage(jid string, user string) (int64, error) {
	result := db.Unscoped().Where("jid = ? AND user = ?", jid, user).Delete(&RepeatedMessage{})
	return result.RowsAffected, result.Error
}

func UpdateNextDate(id uint, nextDate time.Time) error {
	return db.Model(&RepeatedMessage{}).Where("id = ?", id).Update("next_date", nextDate).Error
}

func GetRepeatedMessageToday() []RepeatedMessage {
	var messages []RepeatedMessage

	now := time.Now()
	db.Where("next_date <= ?", now).Find(&messages)

	return messages
}
