package db

import "gorm.io/gorm"

type MessageEvent struct {
	gorm.Model
	JID  string `gorm:"not null;column:jid"`
	Type string `gorm:"not null"`
}

func CreateMessageEvent(jid string, msgType string) {
	db.Create(&MessageEvent{
		JID:  jid,
		Type: msgType,
	})
}

func DeleteMessageEvent(jid string, msgType string) {
	db.Unscoped().Delete(&MessageEvent{}, "jid = ? AND type = ?", jid, msgType)
}

func ToggleMessageEvent(jid string, msgType string) (bool) {
	var msgEvent MessageEvent
	db.Where("jid = ? AND type = ?", jid, msgType).First(&msgEvent)
	if msgEvent.JID == "" {
		CreateMessageEvent(jid, msgType)
		return true
	} else {
		DeleteMessageEvent(jid, msgType)
		return false
	}
}

func GetMessageEvents(jid string) []MessageEvent {
	var msgEvents []MessageEvent
	db.Where("jid = ?", jid).Find(&msgEvents)
	return msgEvents
}
