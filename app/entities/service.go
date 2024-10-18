package entities

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type entitiesApi struct {
	db *gorm.DB
	logger log.AllLogger
}

type EntitiesAPI interface {
	Create(req *CreateEntityRequest) (res *EntityResponse, err error)
	GetEntity(req *GetEntityRequest) (res *EntityResponse, err error)
}

func NewEntitiesAPI(db *gorm.DB, logger log.AllLogger) EntitiesAPI {
	return &entitiesApi{
		db: db,
		logger: logger,
	}
}

// @Summary      	Create Entity
// @Description	Validates name, type. Creates a new entity.
// @Tags			Entities
// @Accept			json
// @Produce			json
// @Param			CreateEntityRequest	body		CreateEntityRequest	true	"CreateEntityRequest"
// @Success			200					{object}	EntityResponse
// @Router			/api/entities/	[POST]
func (s *entitiesApi) Create(req *CreateEntityRequest) (res *EntityResponse, err error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if req.Type == "" {
		return nil, fmt.Errorf("type is required")
	}

	entity := &Entity{
		Name: req.Name,
		Type: req.Type,
	}

	result := s.db.Create(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &EntityResponse{
		ID: entity.ID,
		Name: entity.Name,
		Type: entity.Type,
	}, nil
}

// @Summary      	Get Entity
// @Description	Validates name, type. Creates a new entity.
// @Tags			Entities
// @Accept			json
// @Produce			json
// @Param           id   path int false "ID"
// @Param			name	query		string	false	"Name"
// @Success			200					{object}	EntityResponse
// @Router			/api/entities/{id}	[GET]
func (s *entitiesApi) GetEntity(req *GetEntityRequest) (res *EntityResponse, err error) {
	if req.ID == 0 && req.Name == "" {
		return nil, fmt.Errorf("id or name is required")
	}
	
	entity := &Entity{}

	result := s.db.Where("id = ? OR name = ?", req.ID, req.Name).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &EntityResponse{
		ID: entity.ID,
		Name: entity.Name,
		Type: entity.Type,
	}, nil
}	
