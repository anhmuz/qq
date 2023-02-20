package qqserver

import (
	"qq/models"
	"qq/pkg/protocol"
)

func FromAddMessage(message protocol.AddMessage) models.Entity {
	return models.Entity{
		Key:   message.Key,
		Value: message.Value,
	}
}

func ToAddReplyMessage(added bool) protocol.AddReplyMessage {
	return protocol.AddReplyMessage{
		BaseReplyMessage: protocol.BaseReplyMessage{Name: protocol.AddMessageName},
		Added:            added,
	}
}

func FromRemoveMessage(message protocol.RemoveMessage) string {
	return message.Key
}

func ToRemoveReplyMessage(removed bool) protocol.RemoveReplyMessage {
	return protocol.RemoveReplyMessage{
		BaseReplyMessage: protocol.BaseReplyMessage{Name: protocol.RemoveMessageName},
		Removed:          removed,
	}
}

func FromGetMessage(message protocol.GetMessage) string {
	return message.Key
}

func ToGetReplyMessage(entity *models.Entity) protocol.GetReplyMessage {
	if entity != nil {
		return protocol.GetReplyMessage{
			BaseReplyMessage: protocol.BaseReplyMessage{Name: protocol.GetMessageName},
			Value:            &entity.Value,
		}
	}
	return protocol.GetReplyMessage{
		BaseReplyMessage: protocol.BaseReplyMessage{Name: protocol.GetMessageName},
		Value:            nil,
	}
}

func ToGetAllReplyMessage(entities []models.Entity) protocol.GetAllReplyMessage {
	data := make([]protocol.Entity, 0, len(entities))

	for _, entity := range entities {
		data = append(data, protocol.Entity{
			Key:   entity.Key,
			Value: entity.Value,
		})
	}

	return protocol.GetAllReplyMessage{
		BaseReplyMessage: protocol.BaseReplyMessage{Name: protocol.GetAllMessageName},
		Entities:         data,
	}
}
