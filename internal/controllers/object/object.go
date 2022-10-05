package object

import (
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"github.com/shawnkhoffman/go-modular-rest-api-dynamodb/internal/entities/object"
	"github.com/shawnkhoffman/go-modular-rest-api-dynamodb/internal/repository/adapter"
)

type Controller struct {
	repository adapter.Interface
}

type Interface interface {
	DescribeOne(ID uuid.UUID) (entity object.Object, err error)
	DescribeAll() (entities []object.Object, err error)
	Create(entity *object.Object) (uuid.UUID, error)
	Update(ID uuid.UUID, entity *object.Object) error
	Remove(ID uuid.UUID) error
}

func NewController(repository adapter.Interface) Interface {
	return &Controller{repository: repository}
}

func (c *Controller) DescribeOne(id uuid.UUID) (entity object.Object, err error) {
	entity.ID = id
	response, err := c.repository.FindOne(entity.GetFilterId(), entity.TableName())
	if err != nil {
		return entity, err
	}
	return object.ParseDynamoAtributeToStruct(response.Item)
}

func (c *Controller) DescribeAll() (entities []object.Object, err error) {
	entities = []object.Object{}
	var entity object.Object

	filter := expression.Name("name").NotEqual(expression.Value(""))
	condition, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return entities, err
	}

	response, err := c.repository.FindAll(condition, entity.TableName())
	if err != nil {
		return entities, err
	}

	if response != nil {
		for _, value := range response.Items {
			entity, err := object.ParseDynamoAtributeToStruct(value)
			if err != nil {
				return entities, err
			}
			entities = append(entities, entity)
		}
	}

	return entities, nil
}

func (c *Controller) Create(entity *object.Object) (uuid.UUID, error) {
	entity.CreatedAt = time.Now()
	_, err := c.repository.CreateOrUpdate(entity.GetMap(), entity.TableName())
	return entity.ID, err
}

func (c *Controller) Update(id uuid.UUID, entity *object.Object) error {
	found, err := c.DescribeOne(id)
	if err != nil {
		return err
	}
	found.ID = id
	found.Name = entity.Name
	found.UpdatedAt = time.Now()
	_, err = c.repository.CreateOrUpdate(found.GetMap(), entity.TableName())
	return err
}

func (c *Controller) Remove(id uuid.UUID) error {
	entity, err := c.DescribeOne(id)
	if err != nil {
		return err
	}
	_, err = c.repository.Delete(entity.GetFilterId(), entity.TableName())
	return err
}
