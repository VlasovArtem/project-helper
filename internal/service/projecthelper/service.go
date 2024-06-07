package projecthelper

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks -source=service.go

import (
	"context"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"project-helper/internal/config"
	"project-helper/internal/domain/entity"
)

type (
	ArgService interface {
		PrepareArgs(ctx context.Context, operation config.Operation) ([]string, error)
	}
	OperationService interface {
		GetEnhancedOperation(ctx context.Context, name string) (config.Operation, error)
		GetOperationExecutionPath(ctx context.Context, name string) (string, error)
	}
	FlagService interface {
		GetInitialFlags() *entity.Flags
	}
)

type Service struct {
	operationService OperationService
	flagService      FlagService
	argService       ArgService
}

func NewService(operationService OperationService, flagService FlagService, argService ArgService) *Service {
	return &Service{
		operationService: operationService,
		flagService:      flagService,
		argService:       argService,
	}
}

func (s *Service) Run(ctx context.Context) error {
	flags := s.flagService.GetInitialFlags()

	enhancedOperation, err := s.operationService.GetEnhancedOperation(ctx, flags.Operation)
	if err != nil {
		return errors.Wrap(err, "failed to get enhanced operation")
	}

	return s.runOperation(ctx, enhancedOperation)
}

func (s *Service) runOperation(ctx context.Context, operation config.Operation) error {
	err := s.runBefore(ctx, operation)
	if err != nil {
		return errors.Wrap(err, "failed to run before")
	}

	args, err := s.argService.PrepareArgs(ctx, operation)
	if err != nil {
		return errors.Wrap(err, "failed to prepare args")
	}

	log.Debug().
		Str("operation.description", operation.Description).
		Strs("operation.args.raw", args).
		Str("operation.cmd", operation.Cmd).
		Any("operation.args.predefined", operation.PredefinedFlags).
		Msgf("Running operation")

	if err = s.runCmd(ctx, operation, args); err != nil {
		return errors.Wrap(err, "failed to run command")
	}

	return nil
}

func (s *Service) runBefore(ctx context.Context, operation config.Operation) error {
	for _, runBeforeOperation := range operation.RunBefore {
		err := s.runOperation(ctx, runBeforeOperation)
		if err != nil {
			return errors.Wrapf(err, "failed to run before operation: %s", runBeforeOperation.Name)
		}
	}

	return nil
}

func (s *Service) runCmd(ctx context.Context, operation config.Operation, finalArgs []string) error {
	command := exec.CommandContext(ctx, operation.Cmd, finalArgs...)
	if operation.ChangePath {
		executionPath, err := s.operationService.GetOperationExecutionPath(ctx, operation.Name)
		if err != nil {
			return errors.Wrap(err, "failed to get operation execution path")
		}

		command.Dir = executionPath
	}

	log.Debug().Msgf("Command execution: %s", command.String())

	command.Env = append(os.Environ())

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	if err := command.Run(); err != nil {
		return errors.Wrap(err, "failed to run command")
	}
	return nil
}
