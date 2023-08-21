package db

import (
	"gorm.io/gorm"
)

type MessageEvent struct {
	gorm.Model
	JID     string `gorm:"not null;column:jid;uniqueIndex:jid_content_idx"`
	Content string `gorm:"not null;uniqueIndex:jid_content_idx"`
}

func CreateMessageEvent(jid string, content string) (*MessageEvent, error) {
	data := &MessageEvent{
		JID:     jid,
		Content: content,
	}

	err := db.Create(data).Error

	return data, err
}

func ClearMessageEvents(jid string) error {
	err := db.Unscoped().Delete(&MessageEvent{}, "jid = ?", jid).Error
	return err
}

func GetMessageEvents(jid string) []MessageEvent {
	var msgEvents []MessageEvent
	db.Where("jid = ?", jid).Find(&msgEvents)
	return msgEvents
}
