package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"project-helper/internal/config"
	domainerrors "project-helper/internal/domain/errors"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	applicationConfig := createApplication()

	tests := map[string]struct {
		preconditions func(t *testing.T)
		output        *Service
		expectedError error
	}{
		"with config path variable": {
			preconditions: func(t *testing.T) {
				dir := t.TempDir()
				create, err := os.Create(filepath.Join(dir, "config.yaml"))
				require.NoError(t, err)

				err = yaml.NewEncoder(create).Encode(applicationConfig)

				require.NoError(t, err)

				_ = os.Setenv("CONFIG_PATH", create.Name())
			},
			output: &Service{
				config: &applicationConfig,
				predefinedArgs: map[string]config.PredefinedArg{"predefined-arg-name": {
					Name: "predefined-arg-name", Type: "string", Args: config.Args{{Name: "arg", Values: []string{"value1", "value2"}}}},
				},
				operationsMap: map[string]config.Operation{
					"operation-name": {
						Description:       "description",
						Name:              "operation-name",
						ShortName:         "on",
						Cmd:               "cmd",
						Args:              []string{"arg1", "arg2"},
						ExecutionPath:     "operation-path",
						ChangePath:        false,
						PredefinedArgsTag: &config.PredefinedArgsTag{Name: "predefined-args-tag-name", Value: "predefined-args-tag-value"},
						RunBefore: config.Operations{
							{
								Description:     "run before description",
								Args:            []string{},
								RunBefore:       config.Operations{},
								PredefinedFlags: config.PredefinedFlags{},
							},
						},
						PredefinedFlags: config.PredefinedFlags{{Name: "predefined-flag-name", Value: "predefined-flag-value"}},
					},
					"on": {
						Description:       "description",
						Name:              "operation-name",
						ShortName:         "on",
						Cmd:               "cmd",
						Args:              []string{"arg1", "arg2"},
						ExecutionPath:     "operation-path",
						ChangePath:        false,
						PredefinedArgsTag: &config.PredefinedArgsTag{Name: "predefined-args-tag-name", Value: "predefined-args-tag-value"},
						RunBefore: config.Operations{
							{
								Description:     "run before description",
								Args:            []string{},
								RunBefore:       config.Operations{},
								PredefinedFlags: config.PredefinedFlags{},
							},
						},
						PredefinedFlags: config.PredefinedFlags{{Name: "predefined-flag-name", Value: "predefined-flag-value"}},
					},
				},
				additionalArgs: map[string]string{"application-path": applicationConfig.Path},
			},
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			testCase.preconditions(t)

			svc, err := NewService()

			if testCase.expectedError != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.output.GetConfig(), svc.GetConfig())
				assert.Equal(t, testCase.output.GetPredefinedArgs(), svc.GetPredefinedArgs())
				assert.Equal(t, testCase.output.GetApplicationPath(), svc.GetApplicationPath())
				assert.Equal(t, testCase.output.GetAdditionalArgs(), svc.GetAdditionalArgs())
			}
		})
	}
}

func TestGetOperation(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		service     *Service
		name        string
		expected    config.Operation
		expectedErr error
	}{
		"operation found": {
			service: &Service{
				operationsMap: map[string]config.Operation{
					"operation-name": {
						Description: "description",
					},
				},
			},
			name: "operation-name",
			expected: config.Operation{
				Description: "description",
			},
		},
		"operation not found": {
			service: &Service{
				operationsMap: map[string]config.Operation{},
			},
			name:        "operation-name",
			expectedErr: domainerrors.ErrorOperationNotFound,
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			operation, err := testCase.service.GetOperation(context.Background(), testCase.name)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expected, operation)
			}
		})
	}
}

func createApplication() config.Application {
	return config.Application{
		Name: "test",
		Operations: config.Operations{
			{
				Description:   "description",
				Name:          "operation-name",
				ShortName:     "on",
				Cmd:           "cmd",
				Args:          []string{"arg1", "arg2"},
				ExecutionPath: "operation-path",
				ChangePath:    false,
				PredefinedArgsTag: &config.PredefinedArgsTag{
					Name:  "predefined-args-tag-name",
					Value: "predefined-args-tag-value",
				},
				RunBefore: config.Operations{
					{
						Description:     "run before description",
						Args:            []string{},
						RunBefore:       config.Operations{},
						PredefinedFlags: config.PredefinedFlags{},
					},
				},
				PredefinedFlags: config.PredefinedFlags{
					{
						Name:  "predefined-flag-name",
						Value: "predefined-flag-value",
					},
				},
			},
		},
		Path: "path",
		DynamicFlags: config.DynamicFlags{
			{
				Name:        "dynamic-flag-name",
				ShortName:   "d",
				Description: "dynamic flag description",
				Type:        "string",
				Default:     "true",
			},
		},
		PredefinedArgs: config.PredefinedArgs{
			{
				Name: "predefined-arg-name",
				Type: "string",
				Args: config.Args{
					{
						Name:   "arg",
						Values: []string{"value1", "value2"},
					},
				},
			},
		},
	}
}
