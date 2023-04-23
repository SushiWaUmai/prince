package db

import (
	"time"

	"gorm.io/gorm"
)

type RepeatedMessage struct {
	gorm.Model
	JID      string    `gorm:"not null;column:jid"`
	Message  string    `gorm:"not null"`
	Repeat   string    `gorm:"not null"`
	NextDate time.Time `gorm:"not null"`
}

func CreateRepeatedMessage(jid string, message string, repeat string, nextDate time.Time) {
	db.Create(&RepeatedMessage{
		JID:      jid,
		Message:  message,
		Repeat:   repeat,
		NextDate: nextDate,
	})
}

func ClearRepeatedMessage(jid string) int64 {
	result := db.Delete(&RepeatedMessage{}, "jid = ?", jid)
	return result.RowsAffected
}

func UpdateNextDate(id uint, nextDate time.Time) {
	db.Model(&RepeatedMessage{}).Where("id = ?", id).Update("next_date", nextDate)
}

func GetRepeatedMessageToday() []RepeatedMessage {
	var messages []RepeatedMessage
	now := time.Now()
	today := now.Format("2006-01-02") // format today's date as YYYY-MM-DD
	db.Where("next_date <= ? AND repeat LIKE %?%", now, today).Find(&messages)
	return messages
}
