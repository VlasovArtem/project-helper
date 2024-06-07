package predefined

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks -source=service.go

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"project-helper/internal/config"
	"project-helper/internal/domain/dto"
	domainerrors "project-helper/internal/domain/errors"
	"project-helper/internal/utils"
)

type (
	ConfigService interface {
		GetPredefinedArgs() map[string]config.PredefinedArg
	}
)

type Service struct {
	configService ConfigService
}

func NewService(configService ConfigService) *Service {
	return &Service{
		configService: configService,
	}
}

func (s *Service) TryToFindPredefinedArgValue(request *dto.TryToFindPredefinedArgRequest) (string, error) {
	err := utils.Validate.Struct(request)
	if err != nil {
		return "", errors.Wrap(err, "failed to validate request")
	}

	arg := s.configService.GetPredefinedArgs()[request.ParsedTag]

	predefinedValue, err := arg.Args.GetArgValues(request.Value)
	if err != nil {
		log.Debug().
			Str("tag", request.ParsedTag).
			Str("value", request.Value).
			Err(err).Msgf("Failed to get predefined arg values")

		return request.Value, nil
	}

	return strings.Join(predefinedValue, ","), nil
}

func (s *Service) GetPredefinedArgValues(request *dto.GetPredefinedArgsRequest) ([]string, error) {
	err := utils.Validate.Struct(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate request")
	}

	predefinedArg, ok := s.configService.GetPredefinedArgs()[request.PredefinedArgsTag.Value]
	if !ok {
		return nil, errors.Wrapf(domainerrors.ErrorPredefinedArgNotFound, "predefined arg %s not found", request.PredefinedArgsTag.Value)
	}

	value, err := request.Flags.GetRequiredFlagStringValue(request.PredefinedArgsTag.Name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get flag value")
	}

	values, err := predefinedArg.Args.GetArgValues(value)
	if err != nil {
		values, err = predefinedArg.Args.GetArgValues("*")
		if err == nil {
			return values, nil
		}
		return nil, errors.Wrapf(err, "failed to get arg values for value %s or common value (*)", value)
	}

	return values, nil
}
