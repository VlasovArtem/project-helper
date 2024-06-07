package arg

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"project-helper/internal/config"
	"project-helper/internal/domain/dto"
	"project-helper/internal/domain/entity"
	"project-helper/internal/service/arg/mocks"
	"project-helper/internal/utils"
)

func TestPrepareArgs(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		preconditions func(*testController)
		operation     config.Operation
		expected      []string
		expectedErr   error
	}{
		"valid operation with enhanced operation args": {
			preconditions: func(t *testController) {
				t.flagService.EXPECT().GetOperationFlags(config.Operation{
					Name:            "test",
					PredefinedFlags: config.PredefinedFlags{},
					Args:            []string{"{{predefined_arg1}}", "arg2"},
					PredefinedArgsTag: &config.PredefinedArgsTag{
						Name:  "predefined_arg1",
						Value: "predefined_arg1_value",
					},
				}).
					Return(&entity.Flags{
						DynamicFlags: map[string]any{
							"predefined_arg1": utils.MakePointer("predefined_arg1_name"),
						},
					})
				t.enhanceArgService.EXPECT().GetEnhancedOperationArgs(&dto.GetEnhancedOperationArgs{
					Flags: &entity.Flags{
						DynamicFlags: map[string]any{
							"predefined_arg1": utils.MakePointer("predefined_arg1_name"),
						},
					},
					Operation: config.Operation{
						Name:            "test",
						PredefinedFlags: config.PredefinedFlags{},
						Args:            []string{"{{predefined_arg1}}", "arg2"},
						PredefinedArgsTag: &config.PredefinedArgsTag{
							Name:  "predefined_arg1",
							Value: "predefined_arg1_value",
						},
					},
				}).Return([]string{"enhanced_operation_arg1", "enhanced_operation_arg2"}, nil)
				t.enhanceArgService.EXPECT().EnhanceArgs(&dto.EnhanceArgsRequest{
					Flags: &entity.Flags{
						DynamicFlags: map[string]any{
							"predefined_arg1": utils.MakePointer("predefined_arg1_name"),
						},
					},
					Operation: config.Operation{
						Name:            "test",
						PredefinedFlags: config.PredefinedFlags{},
						Args:            []string{"{{predefined_arg1}}", "arg2"},
						PredefinedArgsTag: &config.PredefinedArgsTag{
							Name:  "predefined_arg1",
							Value: "predefined_arg1_value",
						},
					},
					Args: []string{"enhanced_operation_arg1", "enhanced_operation_arg2"},
				}).Return([]string{"processed_enhanced_operation_arg1", "processed_enhanced_operation_arg2"}, nil)
			},
			operation: config.Operation{
				Name:            "test",
				PredefinedFlags: config.PredefinedFlags{},
				Args:            []string{"{{predefined_arg1}}", "arg2"},
				PredefinedArgsTag: &config.PredefinedArgsTag{
					Name:  "predefined_arg1",
					Value: "predefined_arg1_value",
				},
			},
			expected: []string{"processed_enhanced_operation_arg1", "processed_enhanced_operation_arg2"},
		},
		"valid operation without enhanced operation args": {
			preconditions: func(t *testController) {
				t.flagService.EXPECT().GetOperationFlags(config.Operation{
					Name:            "test",
					PredefinedFlags: config.PredefinedFlags{},
					PredefinedArgsTag: &config.PredefinedArgsTag{
						Name:  "predefined_arg1",
						Value: "predefined_arg1_value",
					},
				}).
					Return(&entity.Flags{
						DynamicFlags: map[string]any{
							"predefined_arg1": utils.MakePointer("predefined_arg1_name"),
						},
					})
				t.predefinedArgSvc.EXPECT().GetPredefinedArgValues(&dto.GetPredefinedArgsRequest{
					Flags: &entity.Flags{
						DynamicFlags: map[string]any{
							"predefined_arg1": utils.MakePointer("predefined_arg1_name"),
						},
					},
					PredefinedArgsTag: &config.PredefinedArgsTag{
						Name:  "predefined_arg1",
						Value: "predefined_arg1_value",
					},
				}).
					Return([]string{"predefined_arg1"}, nil)
				t.enhanceArgService.EXPECT().EnhanceArgs(&dto.EnhanceArgsRequest{
					Flags: &entity.Flags{
						DynamicFlags: map[string]any{
							"predefined_arg1": utils.MakePointer("predefined_arg1_name"),
						},
					},
					Operation: config.Operation{
						Name:            "test",
						PredefinedFlags: config.PredefinedFlags{},
						PredefinedArgsTag: &config.PredefinedArgsTag{
							Name:  "predefined_arg1",
							Value: "predefined_arg1_value",
						},
					},
					Args: []string{"predefined_arg1"},
				}).Return([]string{"enhanced_predefined_arg1"}, nil)
			},
			operation: config.Operation{
				Name:            "test",
				PredefinedFlags: config.PredefinedFlags{},
				PredefinedArgsTag: &config.PredefinedArgsTag{
					Name:  "predefined_arg1",
					Value: "predefined_arg1_value",
				},
			},
			expected: []string{"enhanced_predefined_arg1"},
		},
		"valid operation without predefined args": {
			preconditions: func(t *testController) {
				t.flagService.EXPECT().GetOperationFlags(config.Operation{
					Name:            "test",
					PredefinedFlags: config.PredefinedFlags{},
					Args:            []string{"arg1"},
				}).
					Return(&entity.Flags{
						DynamicFlags: map[string]any{
							"arg1": utils.MakePointer("arg1_name"),
						},
					})
				t.enhanceArgService.EXPECT().EnhanceArgs(&dto.EnhanceArgsRequest{
					Flags: &entity.Flags{
						DynamicFlags: map[string]any{
							"arg1": utils.MakePointer("arg1_name"),
						},
					},
					Operation: config.Operation{
						Name:            "test",
						PredefinedFlags: config.PredefinedFlags{},
						Args:            []string{"arg1"},
					},
					Args: []string{"arg1"},
				}).Return([]string{"enhanced_arg1"}, nil)
			},
			operation: config.Operation{
				Name:            "test",
				PredefinedFlags: config.PredefinedFlags{},
				Args:            []string{"arg1"},
			},
			expected: []string{"enhanced_arg1"},
		},
		"valid operation no args": {
			preconditions: func(t *testController) {
				t.flagService.EXPECT().GetOperationFlags(config.Operation{
					Name:            "test",
					PredefinedFlags: config.PredefinedFlags{},
				}).
					Return(&entity.Flags{
						DynamicFlags: map[string]any{
							"arg1": utils.MakePointer("arg1_name"),
						},
					})
			},
			operation: config.Operation{
				Name:            "test",
				PredefinedFlags: config.PredefinedFlags{},
			},
			expected: []string{},
		},
		"with error on get predefined args": {
			preconditions: func(t *testController) {
				t.flagService.EXPECT().GetOperationFlags(gomock.Any()).
					Return(&entity.Flags{
						DynamicFlags: map[string]any{
							"predefined_arg1": utils.MakePointer("predefined_arg1_name"),
						},
					})
				t.predefinedArgSvc.EXPECT().GetPredefinedArgValues(gomock.Any()).
					Return([]string{}, assert.AnError)
			},
			operation: config.Operation{
				Name:            "test",
				PredefinedFlags: config.PredefinedFlags{},
				PredefinedArgsTag: &config.PredefinedArgsTag{
					Name:  "predefined_arg1",
					Value: "predefined_arg1_value",
				},
			},
			expectedErr: errors.New("failed to enhance args: failed to get predefined args: assert.AnError general error for testing"),
		},
		"with error on get enhanced operation args": {
			preconditions: func(t *testController) {
				t.flagService.EXPECT().GetOperationFlags(gomock.Any()).
					Return(&entity.Flags{
						DynamicFlags: map[string]any{
							"predefined_arg1": utils.MakePointer("predefined_arg1_name"),
						},
					})
				t.enhanceArgService.EXPECT().GetEnhancedOperationArgs(gomock.Any()).
					Return([]string{}, assert.AnError)
			},
			operation: config.Operation{
				Name:            "test",
				PredefinedFlags: config.PredefinedFlags{},
				Args:            []string{"{{predefined_arg1}}", "arg2"},
				PredefinedArgsTag: &config.PredefinedArgsTag{
					Name:  "predefined_arg1",
					Value: "predefined_arg1_value",
				},
			},
			expectedErr: errors.New("failed to enhance args: failed to enhance with operation args: assert.AnError general error for testing"),
		},
		"with error on enhance args": {
			preconditions: func(t *testController) {
				t.flagService.EXPECT().GetOperationFlags(gomock.Any()).
					Return(&entity.Flags{
						DynamicFlags: map[string]any{
							"predefined_arg1": utils.MakePointer("predefined_arg1_name"),
						},
					})
				t.predefinedArgSvc.EXPECT().GetPredefinedArgValues(gomock.Any()).
					Return([]string{"predefined_arg1"}, nil)
				t.enhanceArgService.EXPECT().EnhanceArgs(gomock.Any()).
					Return([]string{}, assert.AnError)
			},
			operation: config.Operation{
				Name:            "test",
				PredefinedFlags: config.PredefinedFlags{},
				PredefinedArgsTag: &config.PredefinedArgsTag{
					Name:  "predefined_arg1",
					Value: "predefined_arg1_value",
				},
			},
			expectedErr: errors.New("failed to enhance args: assert.AnError general error for testing"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			controller := newTestController(ctrl)

			testCase.preconditions(controller)

			service := controller.Build()

			args, err := service.PrepareArgs(nil, testCase.operation)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expected, args)
			}
		})
	}
}

type testController struct {
	flagService       *mocks.MockFlagService
	enhanceArgService *mocks.MockEnhanceArgService
	predefinedArgSvc  *mocks.MockPredefinedArgService
}

func newTestController(ctrl *gomock.Controller) *testController {
	return &testController{
		flagService:       mocks.NewMockFlagService(ctrl),
		enhanceArgService: mocks.NewMockEnhanceArgService(ctrl),
		predefinedArgSvc:  mocks.NewMockPredefinedArgService(ctrl),
	}
}

func (t *testController) Build() *Service {
	return NewService(
		t.flagService,
		t.enhanceArgService,
		t.predefinedArgSvc,
	)
}
