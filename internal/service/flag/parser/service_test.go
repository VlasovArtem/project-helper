package parser

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"project-helper/internal/config"
	"project-helper/internal/domain/entity"
	"project-helper/internal/service/flag/parser/mocks"
	"project-helper/internal/utils"
)

func TestParseFlags(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		precondition  func(*testController)
		args          []string
		expectedFlags *entity.Flags
		expectedError error
	}{
		"valid string flags": {
			precondition: func(t *testController) {
				t.configService.EXPECT().GetConfig().Return(&config.Application{
					DynamicFlags: []config.DynamicFlag{
						{Name: "flag", Type: entity.String},
					},
				})
			},
			args: []string{"--operation=test", "--flag=value"},
			expectedFlags: &entity.Flags{
				Operation:    "test",
				DynamicFlags: map[string]any{"flag": utils.MakePointer("value")},
			},
		},
		"valid array flags": {
			precondition: func(t *testController) {
				t.configService.EXPECT().GetConfig().Return(&config.Application{
					DynamicFlags: []config.DynamicFlag{
						{Name: "flag", Type: entity.Array},
					},
				})
			},
			args: []string{"--operation=test", "--flag=value1,value2"},
			expectedFlags: &entity.Flags{
				Operation:    "test",
				DynamicFlags: map[string]any{"flag": utils.MakePointer([]string{"value1", "value2"})},
			},
		},
		"with missing operation": {
			precondition: func(t *testController) {
				t.configService.EXPECT().GetConfig().Return(&config.Application{
					DynamicFlags: []config.DynamicFlag{},
				})
			},
			expectedFlags: &entity.Flags{
				Operation:    "test",
				DynamicFlags: map[string]any{"flag": utils.MakePointer("value")},
			},
			expectedError: errors.New("failed to validate flags: operation not provided"),
		},
		"unknown flag type": {
			args: []string{"--operation=test"},
			precondition: func(t *testController) {
				t.configService.EXPECT().GetConfig().Return(&config.Application{
					DynamicFlags: []config.DynamicFlag{
						{Name: "flag", Type: "unknown"},
					},
				})
			},
			expectedError: errors.New("unknown flag type unknown"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			flagSet = func() *pflag.FlagSet {
				flags := pflag.NewFlagSet("run", pflag.ContinueOnError)

				flags.ParseErrorsWhitelist.UnknownFlags = true

				return flags
			}()

			controller := newTestController(gomock.NewController(t))

			if testCase.precondition != nil {
				testCase.precondition(controller)
			}

			os.Args = append([]string{os.Args[0]}, testCase.args...)

			service := controller.Build()

			flags, err := service.ParseFlags()

			if testCase.expectedError != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedFlags, flags)
			}
		})
	}
}

type testController struct {
	configService *mocks.MockConfigService
}

func newTestController(ctrl *gomock.Controller) *testController {
	return &testController{
		configService: mocks.NewMockConfigService(ctrl),
	}
}

func (t *testController) Build() *Service {
	return NewService(t.configService)
}
