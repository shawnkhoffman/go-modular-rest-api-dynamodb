package object

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/shawnkhoffman/go-modular-rest-api-dynamodb/internal/entities"
)

type Object struct {
	entities.Base
	Name string `json:"name"`
}

func InterfaceToModel(data interface{}) (instance *Object, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, &instance)
}

func (p *Object) GetFilterId() map[string]interface{} {
	return map[string]interface{}{"_id": p.ID.String()}
}

func (p *Object) TableName() string {
	return "objects"
}

func (p *Object) Bytes() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Object) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"_id":       p.ID.String(),
		"name":      p.Name,
		"createdAt": p.CreatedAt.Format(entities.GetTimeFormat()),
		"updatedAt": p.UpdatedAt.Format(entities.GetTimeFormat()),
	}
}

func ParseDynamoAtributeToStruct(response map[string]*dynamodb.AttributeValue) (p Object, err error) {
	if response == nil || (response != nil && len(response) == 0) {
		return p, errors.New("Item not found")
	}
	for key, value := range response {
		if key == "_id" {
			p.ID, err = uuid.Parse(*value.S)
			if p.ID == uuid.Nil {
				err = errors.New("Item not found")
			}
		}
		if key == "name" {
			p.Name = *value.S
		}
		if key == "createdAt" {
			p.CreatedAt, err = time.Parse(entities.GetTimeFormat(), *value.S)
		}
		if key == "updatedAt" {
			p.UpdatedAt, err = time.Parse(entities.GetTimeFormat(), *value.S)
		}
		if err != nil {
			return p, err
		}
	}

	return p, nil
}
