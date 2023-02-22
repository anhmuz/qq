package http

import (
	"qq/models"
	"qq/pkg/qqclient"
	"qq/pkg/qqclient/http"
)

func FromPostRequest(request http.PostRequest) models.Entity {
	return models.Entity{
		Key:   request.Entity.Key,
		Value: request.Entity.Value,
	}
}

func ToPostResponce(added bool) http.PostResponce {
	return http.PostResponce{
		Added: added,
	}
}

func ToDeleteResponce(removed bool) http.DeleteResponce {
	return http.DeleteResponce{
		Removed: removed,
	}
}

func ToGetResponce(entity *models.Entity) http.GetResponce {
	if entity != nil {
		return http.GetResponce{
			Entity: &qqclient.Entity{Key: entity.Key, Value: entity.Value},
		}
	}
	return http.GetResponce{}
}

func ToGetallResponce(entities []models.Entity) http.GetAllResponce {
	data := make([]qqclient.Entity, 0, len(entities))

	for _, entity := range entities {
		data = append(data, qqclient.Entity{
			Key:   entity.Key,
			Value: entity.Value,
		})
	}

	return http.GetAllResponce{
		Entities: data,
	}
}
