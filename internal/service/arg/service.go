package arg

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks -source=service.go

import (
	"context"

	"github.com/pkg/errors"
	"project-helper/internal/config"
	"project-helper/internal/domain/dto"
	"project-helper/internal/domain/entity"
)

type (
	FlagService interface {
		GetOperationFlags(operation config.Operation) *entity.Flags
	}
	EnhanceArgService interface {
		EnhanceArgs(request *dto.EnhanceArgsRequest) ([]string, error)
		GetEnhancedOperationArgs(request *dto.GetEnhancedOperationArgs) ([]string, error)
	}
	PredefinedArgService interface {
		GetPredefinedArgValues(request *dto.GetPredefinedArgsRequest) ([]string, error)
	}
)

type Service struct {
	flagService       FlagService
	enhanceArgService EnhanceArgService
	predefinedArgSvc  PredefinedArgService
}

func NewService(
	flagService FlagService,
	enhanceArgSvc EnhanceArgService,
	predefinedArgSvc PredefinedArgService,
) *Service {
	return &Service{
		flagService:       flagService,
		enhanceArgService: enhanceArgSvc,
		predefinedArgSvc:  predefinedArgSvc,
	}
}

func (s *Service) PrepareArgs(_ context.Context, operation config.Operation) ([]string, error) {
	flags := s.flagService.GetOperationFlags(operation)

	rawEnhancedArgs, err := s.getArgs(flags, operation)
	if err != nil {
		return nil, errors.Wrap(err, "failed to enhance args")
	}

	if len(rawEnhancedArgs) == 0 {
		return make([]string, 0), nil
	}

	args, err := s.enhanceArgService.EnhanceArgs(&dto.EnhanceArgsRequest{
		Flags:     flags,
		Operation: operation,
		Args:      rawEnhancedArgs,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to enhance args")
	}

	return args, nil
}

func (s *Service) getArgs(flags *entity.Flags, operation config.Operation) ([]string, error) {
	if operation.PredefinedArgsTag == nil {
		return operation.Args, nil
	}

	if len(operation.Args) != 0 {
		if args, err := s.enhanceArgService.GetEnhancedOperationArgs(&dto.GetEnhancedOperationArgs{Flags: flags, Operation: operation}); err != nil {
			return nil, errors.Wrap(err, "failed to enhance with operation args")
		} else {
			return args, nil
		}
	} else {
		if args, err := s.predefinedArgSvc.GetPredefinedArgValues(&dto.GetPredefinedArgsRequest{Flags: flags, PredefinedArgsTag: operation.PredefinedArgsTag}); err != nil {
			return nil, errors.Wrap(err, "failed to get predefined args")
		} else {
			return args, nil
		}
	}
}
