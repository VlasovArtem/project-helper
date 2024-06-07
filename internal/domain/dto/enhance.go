package dto

import (
	"project-helper/internal/config"
	"project-helper/internal/domain/entity"
)

type EnhanceArgsRequest struct {
	Flags     *entity.Flags    `validate:"required"`
	Operation config.Operation `validate:"required"`
	Args      []string         `validate:"required,dive,required"`
}
