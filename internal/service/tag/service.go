package tag

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks -source=service.go

import (
	"path/filepath"

	"github.com/pkg/errors"
	"project-helper/internal/config"
	"project-helper/internal/domain/dto"
	"project-helper/internal/domain/entity"
	domainerrors "project-helper/internal/domain/errors"
	"project-helper/internal/utils"
)

type (
	ConfigService interface {
		GetAdditionalArgs() map[string]string
		GetApplicationPath() string
	}
)

type Service struct {
	configService ConfigService
}

func NewService(
	configService ConfigService,
) *Service {
	return &Service{
		configService: configService,
	}
}

func (s *Service) GetTagValue(request *dto.GetTagValueRequest) (string, error) {
	if err := utils.Validate.Struct(request); err != nil {
		return "", errors.Wrap(err, "request is not valid")
	}

	if flag, err := request.Flags.GetFlag(request.ExtractedTag); err != nil {
		if additionalArg, err := s.checkAdditionalArgs(request.Operation, request.ExtractedTag); err != nil {
			return "", errors.Wrap(err, "failed to check additional args")
		} else {
			return additionalArg, nil
		}
	} else {
		flagStringValue, err := entity.GetString(flag)
		if err != nil {
			return "", errors.Wrap(err, "failed to get flag value")
		}

		return flagStringValue, nil
	}
}

func (s *Service) checkAdditionalArgs(operation config.Operation, tag string) (string, error) {
	additionalArgs := s.configService.GetAdditionalArgs()

	if operation.ChangePath {
		additionalArgs[entity.ExecutionPathTag] = filepath.Join(s.configService.GetApplicationPath(), operation.ExecutionPath)
	}

	if value, ok := additionalArgs[tag]; !ok {
		return "", errors.Wrapf(domainerrors.ErrorAdditionalArgNotFound, "additional arg %s not found", tag)
	} else {
		return value, nil
	}
}
