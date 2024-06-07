package tag

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks -source=service.go

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"project-helper/internal/config"
	"project-helper/internal/domain/dto"
	"project-helper/internal/domain/entity"
	domainerrors "project-helper/internal/domain/errors"
	"project-helper/internal/utils"
)

type (
	ConfigService interface {
		GetPatternTags() map[string]config.PatternTag
		GetAdditionalArgs() map[string]string
		GetApplicationPath() string
	}
)

type Service struct {
	configService    ConfigService
}

func NewService(
	configService ConfigService,
) *Service {
	return &Service{
		configService:    configService,
	}
}

func (s *Service) GetTagValue(request *dto.GetTagValueRequest) (string, error) {
	if err := utils.Validate.Struct(request); err != nil {
		return "", errors.Wrap(err, "request is not valid")
	}

	patternTag, ok := s.configService.GetPatternTags()[request.ExtractedTag]
	if !ok {
		if additionalArg, err := s.checkAdditionalArgs(request.Operation, request.ExtractedTag); err != nil {
			return "", errors.Wrap(err, "failed to check additional args")
		} else {
			return additionalArg, nil
		}
	}

	switch patternTag.Type {
	case entity.String:
		value, err := request.Flags.GetRequiredFlagStringValue(patternTag.Name)
		if err != nil {
			return "", errors.Wrap(err, "failed to get flag value")
		}

		return value, nil
	case entity.Array:
		value, err := request.Flags.GetRequiredFlagArrayValue(patternTag.Name)
		if err != nil {
			return "", errors.Wrap(err, "failed to get flag value")
		}

		return strings.Join(value, patternTag.GetJoin()), nil
	default:
		return "", errors.Wrapf(domainerrors.ErrorPatternTagTypeNotFound, "pattern tag type %s not found", patternTag.Type)
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
