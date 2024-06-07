package flag

import (
	"project-helper/internal/config"
	"project-helper/internal/domain/entity"
)

type Service struct {
	initialFlags   *entity.Flags
	operationFlags map[string]*entity.Flags
}

func NewFlagsService(initialFlags *entity.Flags) *Service {
	return &Service{
		initialFlags:   initialFlags,
		operationFlags: make(map[string]*entity.Flags),
	}
}

func (s *Service) GetInitialFlags() *entity.Flags {
	return s.initialFlags
}

func (s *Service) GetOperationFlags(operation config.Operation) *entity.Flags {
	if s.operationFlags[operation.Name] == nil {
		s.operationFlags[operation.Name] = s.enhanceFlags(operation.PredefinedFlags)
	}
	return s.operationFlags[operation.Name]
}

func (s *Service) enhanceFlags(operationPredefinedFlags config.PredefinedFlags) *entity.Flags {
	if len(operationPredefinedFlags) == 0 {
		return s.initialFlags
	}

	newFlags := *s.initialFlags

	for _, flag := range operationPredefinedFlags {
		newFlags.DynamicFlags[flag.Name] = &flag.Value
	}

	return &newFlags
}
