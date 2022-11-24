package rabbitqq

import (
	"qq/models"
	"qq/pkg/rabbitqq"
)

func FromAddMessage(message rabbitqq.AddMessage) models.Entity {
	return models.Entity{
		Key:   message.Key,
		Value: message.Value,
	}
}

func ToAddReplyMessage(added bool) rabbitqq.AddReplyMessage {
	return rabbitqq.AddReplyMessage{
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: rabbitqq.AddMessageName},
		Added:            added,
	}
}

func FromRemoveMessage(message rabbitqq.RemoveMessage) string {
	return message.Key
}

func ToRemoveReplyMessage(removed bool) rabbitqq.RemoveReplyMessage {
	return rabbitqq.RemoveReplyMessage{
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: rabbitqq.RemoveMessageName},
		Removed:          removed,
	}
}

func FromGetMessage(message rabbitqq.GetMessage) string {
	return message.Key
}

func ToGetReplyMessage(entity *models.Entity) rabbitqq.GetReplyMessage {
	if entity != nil {
		return rabbitqq.GetReplyMessage{
			BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: rabbitqq.GetMessageName},
			Value:            &entity.Value,
		}
	}
	return rabbitqq.GetReplyMessage{
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: rabbitqq.GetMessageName},
		Value:            nil,
	}
}

func ToGetAllReplyMessage(entities []models.Entity) rabbitqq.GetAllReplyMessage {
	data := make([]rabbitqq.Entity, 0, len(entities))

	for _, entity := range entities {
		data = append(data, rabbitqq.Entity{
			Key:   entity.Key,
			Value: entity.Value,
		})
	}

	return rabbitqq.GetAllReplyMessage{
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: rabbitqq.GetAllMessageName},
		Entities:         data,
	}
}
