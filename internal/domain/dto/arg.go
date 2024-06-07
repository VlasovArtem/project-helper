package dto

import (
	"project-helper/internal/config"
	"project-helper/internal/domain/entity"
)

type GetPredefinedArgsRequest struct {
	Flags             *entity.Flags             `validate:"required"`
	PredefinedArgsTag *config.PredefinedArgsTag `validate:"required"`
}

type TryToFindPredefinedArgRequest struct {
	ParsedTag string `validate:"required"`
	Value     string `validate:"required"`
}

type GetEnhancedOperationArgs struct {
	Flags     *entity.Flags    `validate:"required"`
	Operation config.Operation `validate:"required"`
}
