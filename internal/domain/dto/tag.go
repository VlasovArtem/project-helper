package dto

import (
	"project-helper/internal/config"
	"project-helper/internal/domain/entity"
)

type GetTagValueRequest struct {
	Operation    config.Operation `validate:"required"`
	Flags        *entity.Flags    `validate:"required"`
	ExtractedTag string           `validate:"required"`
}
