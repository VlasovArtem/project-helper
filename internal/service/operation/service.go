package operation

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks -source=service.go

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"project-helper/internal/config"
)

type ConfigService interface {
	GetOperation(ctx context.Context, name string) (config.Operation, error)
	GetApplicationPath() string
}

type Service struct {
	configService ConfigService
}

func NewService(configService ConfigService) *Service {
	return &Service{
		configService: configService,
	}
}

func (s *Service) GetEnhancedOperation(ctx context.Context, name string) (config.Operation, error) {
	operation, err := s.configService.GetOperation(ctx, name)
	if err != nil {
		return config.Operation{}, errors.Wrapf(err, "operation %s not found", name)
	}

	var enhancedBeforeOperations config.Operations

	for _, beforeOperation := range operation.RunBefore {
		runBeforeOperation, err := s.configService.GetOperation(ctx, beforeOperation.Name)
		if err != nil {
			return config.Operation{}, errors.Wrapf(err, "before operation %s not found", beforeOperation.Name)
		}
		runBeforeOperation.PredefinedFlags = beforeOperation.PredefinedFlags

		enhancedBeforeOperations = append(enhancedBeforeOperations, runBeforeOperation)
	}

	operation.RunBefore = enhancedBeforeOperations

	return operation, nil
}

func (s *Service) GetOperationExecutionPath(ctx context.Context, name string) (string, error) {
	operation, err := s.configService.GetOperation(ctx, name)
	if err != nil {
		return "", errors.Wrap(err, "failed to get operation")
	}

	executionPath := filepath.Join(s.configService.GetApplicationPath(), operation.ExecutionPath)

	_, err = os.Stat(executionPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to get operation execution path")
	}

	return executionPath, nil
}
