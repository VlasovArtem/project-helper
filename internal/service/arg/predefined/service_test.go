package predefined

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"project-helper/internal/config"
	"project-helper/internal/domain/dto"
	"project-helper/internal/domain/entity"
	"project-helper/internal/service/arg/predefined/mocks"
	"project-helper/internal/utils"
)

func TestTryToFindPredefinedArg(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		preconditions func(*testController)
		request       *dto.TryToFindPredefinedArgRequest
		expected      string
		expectedErr   error
	}{
		"valid request": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetPredefinedArgs().Return(map[string]config.PredefinedArg{
					"parsed_tag": {
						Args: config.Args{
							{
								Name:   "parsed_tag_value",
								Values: []string{"parsed_tag_value1", "parsed_tag_value2"},
							},
						},
					},
				})
			},
			request: &dto.TryToFindPredefinedArgRequest{
				ParsedTag: "parsed_tag",
				Value:     "parsed_tag_value",
			},
			expected: "parsed_tag_value1,parsed_tag_value2",
		},
		"valid request with no predefined args": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetPredefinedArgs().
					Return(map[string]config.PredefinedArg{})
			},
			request: &dto.TryToFindPredefinedArgRequest{
				ParsedTag: "parsed_tag",
				Value:     "parsed_tag_value",
			},
			expected: "parsed_tag_value",
		},
		"invalid request": {
			preconditions: func(t *testController) {},
			request:       &dto.TryToFindPredefinedArgRequest{},
			expectedErr:   errors.New("failed to validate request"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			controller := newTestController(ctrl)
			testCase.preconditions(controller)

			service := controller.Build()

			result, err := service.TryToFindPredefinedArgValue(testCase.request)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, result)
			}
		})
	}
}

func TestGetPredefinedArgs(t *testing.T) {
	t.Parallel()

	flags := &entity.Flags{
		DynamicFlags: map[string]*entity.DynamicFlagValue{
			"predefined_arg1": {
				Name:  "predefined_arg1",
				Type:  entity.String,
				Value: utils.MakePointer("predefined_flag_value"),
			},
		},
	}

	tests := map[string]struct {
		preconditions func(*testController)
		request       *dto.GetPredefinedArgsRequest
		expected      []string
		expectedErr   error
	}{
		"valid request": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetPredefinedArgs().Return(map[string]config.PredefinedArg{
					"predefined_arg1_value": {
						Args: config.Args{
							{
								Name:   "predefined_flag_value",
								Values: []string{"value1", "value2"},
							},
						},
					},
				})
			},
			request: &dto.GetPredefinedArgsRequest{
				PredefinedArgsTag: &config.PredefinedArgsTag{
					Name:  "predefined_arg1",
					Value: "predefined_arg1_value",
				},
				Flags: flags,
			},
			expected: []string{"value1", "value2"},
		},
		"valid request with common args": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetPredefinedArgs().Return(map[string]config.PredefinedArg{
					"predefined_arg1_value": {
						Args: config.Args{
							{
								Name:   "*",
								Values: []string{"wildcard"},
							},
						},
					},
				})
			},
			request: &dto.GetPredefinedArgsRequest{
				PredefinedArgsTag: &config.PredefinedArgsTag{
					Name:  "predefined_arg1",
					Value: "predefined_arg1_value",
				},
				Flags: flags,
			},
			expected: []string{"wildcard"},
		},
		"without common args": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetPredefinedArgs().Return(map[string]config.PredefinedArg{
					"predefined_arg1_value": {
						Args: config.Args{},
					},
				})
			},
			request: &dto.GetPredefinedArgsRequest{
				PredefinedArgsTag: &config.PredefinedArgsTag{
					Name:  "predefined_arg1",
					Value: "predefined_arg1_value",
				},
				Flags: flags,
			},
			expectedErr: errors.New("failed to get arg values for value predefined_flag_value or common value (*): arg * not found"),
		},
		"without required flag": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetPredefinedArgs().Return(map[string]config.PredefinedArg{
					"predefined_arg1_value": {
						Args: config.Args{},
					},
				})
			},
			request: &dto.GetPredefinedArgsRequest{
				PredefinedArgsTag: &config.PredefinedArgsTag{},
				Flags:             flags,
			},
			expectedErr: errors.New("predefined arg  not found: predefined arg not found"),
		},
		"without required flag string value": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetPredefinedArgs().Return(map[string]config.PredefinedArg{
					"predefined_arg1_value": {
						Args: config.Args{
							{
								Name:   "predefined_flag_value",
								Values: []string{"value1", "value2"},
							},
						},
					},
				})
			},
			request: &dto.GetPredefinedArgsRequest{
				PredefinedArgsTag: &config.PredefinedArgsTag{
					Name:  "predefined_arg1",
					Value: "predefined_arg1_value",
				},
				Flags: &entity.Flags{},
			},
			expectedErr: errors.New("failed to get flag value: flag predefined_arg1 not found"),
		},
		"with err on get predefined args": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetPredefinedArgs().Return(map[string]config.PredefinedArg{})
			},
			request: &dto.GetPredefinedArgsRequest{
				PredefinedArgsTag: &config.PredefinedArgsTag{},
				Flags:             &entity.Flags{},
			},
			expectedErr: errors.New("predefined arg  not found: predefined arg not found"),
		},
		"with invalid request": {
			preconditions: func(t *testController) {},
			request:       &dto.GetPredefinedArgsRequest{},
			expectedErr:   errors.New("failed to validate request"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			controller := newTestController(ctrl)
			testCase.preconditions(controller)

			service := controller.Build()

			result, err := service.GetPredefinedArgValues(testCase.request)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, result)
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
