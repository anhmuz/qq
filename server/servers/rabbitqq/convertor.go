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
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: rabbitqq.AddMessageName},
		Added:            added,
	}
}

func (c convertor) RemoveMessageToString(message rabbitqq.RemoveMessage) string {
	return message.Key
}

func (c convertor) BoolToRemoveReplyMessage(removed bool) rabbitqq.RemoveReplyMessage {
	return rabbitqq.RemoveReplyMessage{
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: rabbitqq.RemoveMessageName},
		Removed:          removed,
	}
}

func (c convertor) GetMessageToString(message rabbitqq.GetMessage) string {
	return message.Key
}

func (c convertor) EntityToGetReplyMessage(entity *models.Entity) rabbitqq.GetReplyMessage {
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

func (c convertor) EntitiesToGetAllReplyMessage(entities []models.Entity) rabbitqq.GetAllReplyMessage {
	data := []rabbitqq.EntityItem{}

	for _, entity := range entities {
		data = append(data, rabbitqq.EntityItem{
			Key:   entity.Key,
			Value: entity.Value,
		})
	}

	return rabbitqq.GetAllReplyMessage{
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: rabbitqq.GetAllMessageName},
		Entities:         data,
	}
}
