package parser

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks -source=service.go

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"project-helper/internal/config"
	"project-helper/internal/domain/entity"
)

type ConfigService interface {
	GetConfig() *config.Application
}

var flagSet = func() *pflag.FlagSet {
	flags := pflag.NewFlagSet("run", pflag.ContinueOnError)

	flags.ParseErrorsWhitelist.UnknownFlags = true

	return flags
}()

type Service struct {
	configService ConfigService
}

func NewService(configService ConfigService) *Service {
	return &Service{
		configService: configService,
	}
}

func (s *Service) ParseFlags() (*entity.Flags, error) {
	flags, err := s.parseFlags()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse flags")
	}

	err = flags.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate flags")
	}

	return flags, nil
}

func (s *Service) parseFlags() (*entity.Flags, error) {
	flags := entity.NewFlags()

	flagSet.StringVarP(&flags.Operation, "operation", "o", "", "Operation to run")

	applicationConfig := s.configService.GetConfig()

	for _, dynamicFlag := range applicationConfig.DynamicFlags {
		switch dynamicFlag.Type {
		case entity.String:
			var value string

			flagSet.StringVarP(&value, dynamicFlag.Name, dynamicFlag.ShortName, dynamicFlag.Default, dynamicFlag.Description)

			flags.DynamicFlags[dynamicFlag.Name] = &entity.DynamicFlagValue{
				Value: &value,
				Name:  dynamicFlag.Name,
				Type:  dynamicFlag.Type,
			}
		case entity.Array:
			var value []string

			flagSet.StringSliceVarP(&value, dynamicFlag.Name, dynamicFlag.ShortName, []string{}, dynamicFlag.Description)

			flags.DynamicFlags[dynamicFlag.Name] = &entity.DynamicFlagValue{
				Value: &value,
				Name:  dynamicFlag.Name,
				Type:  dynamicFlag.Type,
			}
		default:
			return nil, errors.Errorf("unknown flag type %s", dynamicFlag.Type)
		}

	}

	err := flagSet.Parse(os.Args[1:])

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse flags")
	}

	return flags, nil
}
