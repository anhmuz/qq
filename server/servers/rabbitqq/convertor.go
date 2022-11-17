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
	EntitiesToGetAllReplyMessage(map[string]models.Entity) rabbitqq.GetAllReplyMessage
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
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: "add reply"},
		Added:            added,
	}
}

func (c convertor) RemoveMessageToString(message rabbitqq.RemoveMessage) string {
	return message.Key
}

func (c convertor) BoolToRemoveReplyMessage(removed bool) rabbitqq.RemoveReplyMessage {
	return rabbitqq.RemoveReplyMessage{
		BaseReplyMessage: rabbitqq.BaseReplyMessage{Name: "remove reply"},
		Removed:          removed,
	}
}

func (c convertor) GetMessageToString(message rabbitqq.GetMessage) string {
	return message.Key
}

func (c convertor) EntityToGetReplyMessage(entity *models.Entity) rabbitqq.GetReplyMessage {
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

func (c convertor) EntitiesToGetAllReplyMessage(entities map[string]models.Entity) rabbitqq.GetAllReplyMessage {
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
