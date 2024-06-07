package operation

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"project-helper/internal/config"
	"project-helper/internal/service/operation/mocks"
)

func TestGetEnhancedOperation(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		preconditions func(*testController)
		name          string
		output        config.Operation
		expectedErr   error
	}{
		"success": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetOperation(gomock.Any(), "operation").
					Return(config.Operation{
						Name: "operation",
						RunBefore: []config.Operation{
							{
								Name: "run-before-operation",
								PredefinedFlags: config.PredefinedFlags{
									{
										Name:  "flag",
										Value: "value",
									},
								},
							},
						},
					}, nil)
				t.configService.EXPECT().GetOperation(gomock.Any(), "run-before-operation").
					Return(config.Operation{Name: "run-before-operation"}, nil)
			},
			name: "operation",
			output: config.Operation{
				Name: "operation",
				RunBefore: []config.Operation{
					{
						Name: "run-before-operation",
						PredefinedFlags: config.PredefinedFlags{
							{
								Name:  "flag",
								Value: "value",
							},
						},
					},
				},
			},
		},
		"success without run before operations": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetOperation(gomock.Any(), "operation").
					Return(config.Operation{Name: "operation"}, nil)
			},
			name: "operation",
			output: config.Operation{
				Name: "operation",
			},
		},
		"with run before operation not found": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetOperation(gomock.Any(), "operation").
					Return(config.Operation{
						Name: "operation",
						RunBefore: []config.Operation{
							{
								Name: "run-before-operation",
								PredefinedFlags: config.PredefinedFlags{
									{
										Name:  "flag",
										Value: "value",
									},
								},
							},
						},
					}, nil)
				t.configService.EXPECT().GetOperation(gomock.Any(), "run-before-operation").
					Return(config.Operation{}, assert.AnError)
			},
			name:        "operation",
			expectedErr: errors.New("before operation run-before-operation not found: assert.AnError general error for testing"),
		},
		"with basic operation not found": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetOperation(gomock.Any(), "operation").
					Return(config.Operation{}, assert.AnError)
			},
			name:        "operation",
			expectedErr: errors.New("operation operation not found: assert.AnError general error for testing"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			controller := newTestController(ctrl)
			testCase.preconditions(controller)

			service := controller.Build()

			operation, err := service.GetEnhancedOperation(context.Background(), testCase.name)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.output, operation)
			}
		})
	}
}

func TestGetOperationExecutionPath(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		preconditions func(*testing.T, *testController)
		name          string
		output        string
		expectedErr   error
	}{
		"success": {
			preconditions: func(t *testing.T, controller *testController) {
				dir := t.TempDir()
				_, err := os.Create(filepath.Join(dir, "operation-path"))
				require.NoError(t, err)

				controller.configService.EXPECT().GetOperation(gomock.Any(), "operation").
					Return(config.Operation{
						Name:          "operation",
						ExecutionPath: "operation-path",
					}, nil)
				controller.configService.EXPECT().GetApplicationPath().
					Return(dir)
			},
			name:   "operation",
			output: "operation-path",
		},
		"with file not found": {
			preconditions: func(t *testing.T, controller *testController) {
				controller.configService.EXPECT().GetOperation(gomock.Any(), "operation").
					Return(config.Operation{
						Name:          "operation",
						ExecutionPath: "operation-path",
					}, nil)
				controller.configService.EXPECT().GetApplicationPath().
					Return("application-path")
			},
			name:        "operation",
			expectedErr: errors.New("failed to get operation execution path: stat application-path/operation-path: no such file or directory"),
		},
		"with operation not found": {
			preconditions: func(t *testing.T, controller *testController) {
				controller.configService.EXPECT().GetOperation(gomock.Any(), "operation").
					Return(config.Operation{}, assert.AnError)
			},
			name:        "operation",
			expectedErr: errors.New("failed to get operation: assert.AnError general error for testing"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			controller := newTestController(ctrl)
			testCase.preconditions(t, controller)

			service := controller.Build()

			executionPath, err := service.GetOperationExecutionPath(context.Background(), testCase.name)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedErr.Error())
			} else {
				require.NoError(t, err)
				assert.Contains(t, executionPath, testCase.output)
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
