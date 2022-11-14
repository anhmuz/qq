package rabbitqq

import (
	"qq/models"
	"qq/pkg/rabbitqq"
)

type Convertor interface {
	AddMessageToEntity(rabbitqq.AddMessage) models.Entity
	BoolToAddReplyMessage(bool) rabbitqq.AddReplyMessage
	RemoveMessageToString(rabbitqq.RemoveMessage) string
	BoolToRemoveReplyMessage(bool) rabbitqq.RemoveReplyMessage
	GetMessageToString(rabbitqq.GetMessage) string
	EntityToGetReplyMessage(*models.Entity) rabbitqq.GetReplyMessage
	EntitiesToGetAllReplyMessage([]models.Entity) rabbitqq.GetAllReplyMessage
}

type convertor struct {
}

var _ Convertor = convertor{}

func NewConvertor() (Convertor, error) {
	return &convertor{}, nil
}

func (c convertor) AddMessageToEntity(message rabbitqq.AddMessage) models.Entity {
	return models.Entity{
		Key:   message.Key,
		Value: message.Value,
	}
}

func (c convertor) BoolToAddReplyMessage(added bool) rabbitqq.AddReplyMessage {
	return rabbitqq.AddReplyMessage{
		Added: added,
	}
}

func (c convertor) RemoveMessageToString(message rabbitqq.RemoveMessage) string {
	return message.Key
}

func (c convertor) BoolToRemoveReplyMessage(removed bool) rabbitqq.RemoveReplyMessage {
	return rabbitqq.RemoveReplyMessage{
		Removed: removed,
	}
}

func (c convertor) GetMessageToString(message rabbitqq.GetMessage) string {
	return message.Key
}

func (c convertor) EntityToGetReplyMessage(entity *models.Entity) rabbitqq.GetReplyMessage {
	if entity != nil {
		return rabbitqq.GetReplyMessage{
			Value: &entity.Value,
		}
	}
	return rabbitqq.GetReplyMessage{
		Value: nil,
	}
}

func (c convertor) EntitiesToGetAllReplyMessage(entities []models.Entity) rabbitqq.GetAllReplyMessage {
	data := make(map[string]string, 0)
	for _, entity := range entities {
		data[entity.Key] = entity.Value
	}
	return rabbitqq.GetAllReplyMessage{
		Entities: data,
	}
}