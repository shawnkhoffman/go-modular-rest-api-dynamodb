package object

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/shawnkhoffman/go-modular-rest-api-dynamodb/internal/controllers/object"
	EntityObject "github.com/shawnkhoffman/go-modular-rest-api-dynamodb/internal/entities/object"
	"github.com/shawnkhoffman/go-modular-rest-api-dynamodb/internal/handlers"
	"github.com/shawnkhoffman/go-modular-rest-api-dynamodb/internal/repository/adapter"
	Rules "github.com/shawnkhoffman/go-modular-rest-api-dynamodb/internal/rules"
	RulesObject "github.com/shawnkhoffman/go-modular-rest-api-dynamodb/internal/rules/object"
	HttpStatus "github.com/shawnkhoffman/go-modular-rest-api-dynamodb/utils/http"
)

type Handler struct {
	handlers.Interface

	Controller object.Interface
	Rules      Rules.Interface
}

func NewHandler(repository adapter.Interface) handlers.Interface {
	return &Handler{
		Controller: object.NewController(repository),
		Rules:      RulesObject.NewRules(),
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if chi.URLParam(r, "ID") != "" {
		h.getOne(w, r)
	} else {
		h.getAll(w, r)
	}
}

func (h *Handler) getOne(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	response, err := h.Controller.DescribeOne(ID)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	response, err := h.Controller.DescribeAll()
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	objectBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	ID, err := h.Controller.Create(objectBody)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, map[string]interface{}{"id": ID.String()})
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	objectBody, err := h.getBodyAndValidate(r, ID)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	if err := h.Controller.Update(ID, objectBody); err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	if err := h.Controller.Remove(ID); err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	HttpStatus.StatusNoContent(w, r)
}

func (h *Handler) getBodyAndValidate(r *http.Request, ID uuid.UUID) (*EntityObject.Object, error) {
	objectBody := &EntityObject.Object{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, objectBody)
	if err != nil {
		return &EntityObject.Object{}, errors.New("body is required")
	}

	objectParsed, err := EntityObject.InterfaceToModel(body)
	if err != nil {
		return &EntityObject.Object{}, errors.New("error on convert body to model")
	}

	setDefaultValues(objectParsed, ID)

	return objectParsed, h.Rules.Validate(objectParsed)
}

func setDefaultValues(object *EntityObject.Object, ID uuid.UUID) {
	object.UpdatedAt = time.Now()
	if ID == uuid.Nil {
		object.ID = uuid.New()
		object.CreatedAt = time.Now()
	} else {
		object.ID = ID
	}
}
