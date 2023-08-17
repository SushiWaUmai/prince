package db

import "gorm.io/gorm"

type MessageEvent struct {
	gorm.Model
	JID  string `gorm:"not null;column:jid"`
	Type string `gorm:"not null"`
}

func CreateMessageEvent(jid string, msgType string) (*MessageEvent, error) {
	data := &MessageEvent{
		JID:  jid,
		Type: msgType,
	}
	err := db.Create(data).Error

	return data, err
}

func DeleteMessageEvent(jid string, msgType string) error {
	err := db.Unscoped().Delete(&MessageEvent{}, "jid = ? AND type = ?", jid, msgType).Error
	return err
}

func ToggleMessageEvent(jid string, msgType string) (bool, error) {
	var msgEvent MessageEvent
	err := db.Where("jid = ? AND type = ?", jid, msgType).First(&msgEvent).Error
	if err != nil {
		return false, err
	}

	if msgEvent.JID == "" {
		_, err := CreateMessageEvent(jid, msgType)
		return true, err
	} else {
		err := DeleteMessageEvent(jid, msgType)
		return false, err
	}
}

func GetMessageEvents(jid string) ([]MessageEvent, error) {
	var msgEvents []MessageEvent
	err := db.Where("jid = ?", jid).Find(&msgEvents).Error
	return msgEvents, err
}
