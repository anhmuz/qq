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
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: "add reply"},
		Added:            added,
	}
}

func FromRemoveMessage(message rabbitqq.RemoveMessage) string {
	return message.Key
}

func ToRemoveReplyMessage(removed bool) rabbitqq.RemoveReplyMessage {
	return rabbitqq.RemoveReplyMessage{
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: "remove reply"},
		Removed:          removed,
	}
}

func FromGetMessage(message rabbitqq.GetMessage) string {
	return message.Key
}

func ToGetReplyMessage(entity *models.Entity) rabbitqq.GetReplyMessage {
	if entity != nil {
		return rabbitqq.GetReplyMessage{
			BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: "get reply"},
			Value:            &entity.Value,
		}
	}
	return rabbitqq.GetReplyMessage{
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: "get reply"},
		Value:            nil,
	}
}

func ToGetAllReplyMessage(entities map[string]models.Entity) rabbitqq.GetAllReplyMessage {
	data := []struct {
		Key   string
		Value string
	}{}
	for _, entity := range entities {
		data = append(data, struct {
			Key   string
			Value string
		}{Key: entity.Key, Value: entity.Value})
	}
	return rabbitqq.GetAllReplyMessage{
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: "get all reply"},
		Entities:         data,
	}
}
