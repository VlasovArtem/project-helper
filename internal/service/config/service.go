package config

import (
	"context"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"project-helper/internal/config"
	"project-helper/internal/domain/entity"
	domainerrors "project-helper/internal/domain/errors"
)

type Service struct {
	config         *config.Application
	patternTags    map[string]config.PatternTag
	predefinedArgs map[string]config.PredefinedArg
	operationsMap  map[string]config.Operation
	additionalArgs map[string]string
}

func NewService() (*Service, error) {
	svc := &Service{}

	err := initService(svc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize service")
	}

	return svc, nil
}

func initService(s *Service) error {
	configPath, err := s.readConfigPath()
	if err != nil {
		return errors.Wrap(err, "failed to read config path")
	}

	file, err := os.Open(configPath)
	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}

	s.config = &config.Application{}

	err = yaml.NewDecoder(file).Decode(s.config)
	if err != nil {
		return errors.Wrap(err, "failed to decode config file")
	}

	s.patternTags = s.config.GetPatternTags()
	s.predefinedArgs = s.config.GetPredefinedArgs()
	s.operationsMap = s.config.GetOperationsMap()
	s.additionalArgs = map[string]string{
		entity.ApplicationPathTag: s.config.Path,
	}

	return nil
}

func (s *Service) GetConfig() *config.Application {
	return s.config
}

func (s *Service) GetPatternTags() map[string]config.PatternTag {
	return s.patternTags
}

func (s *Service) GetPredefinedArgs() map[string]config.PredefinedArg {
	return s.predefinedArgs
}

func (s *Service) GetOperation(_ context.Context, name string) (config.Operation, error) {
	operation, ok := s.operationsMap[name]
	if !ok {
		return config.Operation{}, errors.Wrapf(domainerrors.ErrorOperationNotFound, "operation %s not found", name)
	}

	return operation, nil
}

func (s *Service) GetApplicationPath() string {
	return s.config.Path
}

func (s *Service) GetAdditionalArgs() map[string]string {
	args := make(map[string]string, len(s.additionalArgs))

	for k, v := range s.additionalArgs {
		args[k] = v
	}

	return args
}

func (s *Service) readConfigPath() (string, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath, _ = xdg.SearchConfigFile("project-helper/application.yaml")
	}

	if configPath == "" {
		return "", errors.Errorf("config file not found. 'CONFIG_PATH' environment variable and %s", filepath.Join(xdg.ConfigHome, "project-helper/application.yaml"))
	}

	if _, err := os.Stat(configPath); err != nil {
		return "", errors.Wrapf(err, "config file not found by path %s", configPath)
	}

	log.Debug().Str("config.path", configPath).Msg("config file found")

	return configPath, nil
}
