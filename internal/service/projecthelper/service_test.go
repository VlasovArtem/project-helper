package projecthelper

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"project-helper/internal/config"
	"project-helper/internal/domain/entity"
	"project-helper/internal/service/projecthelper/mocks"
)

func TestRun(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	tests := map[string]struct {
		preconditions func(*testController)
		expectedErr   error
	}{
		"success": {
			preconditions: func(t *testController) {
				t.flagService.EXPECT().GetInitialFlags().
					Return(&entity.Flags{
						Operation: "operation",
					})
				t.operationService.EXPECT().GetEnhancedOperation(gomock.Any(), "operation").
					Return(config.Operation{
						Name: "operation",
						Cmd:  "echo",
						RunBefore: config.Operations{{
							Name:       "before",
							Cmd:        "echo",
							ChangePath: true,
						}},
					}, nil)
				t.argService.EXPECT().PrepareArgs(gomock.Any(), config.Operation{
					Name:       "before",
					Cmd:        "echo",
					ChangePath: true,
				}).Return([]string{"'Hello, World! before'"}, nil)
				t.operationService.EXPECT().GetOperationExecutionPath(gomock.Any(), "before").
					Return(dir, nil)
				t.argService.EXPECT().PrepareArgs(gomock.Any(), config.Operation{
					Name: "operation",
					Cmd:  "echo",
					RunBefore: config.Operations{{
						Name:       "before",
						Cmd:        "echo",
						ChangePath: true,
					}},
				}).Return([]string{"'Hello, World!'"}, nil)
			},
		},
		"with error on get operation execution path": {
			preconditions: func(t *testController) {
				t.flagService.EXPECT().GetInitialFlags().
					Return(&entity.Flags{
						Operation: "operation",
					})
				t.operationService.EXPECT().GetEnhancedOperation(gomock.Any(), "operation").
					Return(config.Operation{
						Name: "operation",
						Cmd:  "echo",
						RunBefore: config.Operations{{
							Name:       "before",
							Cmd:        "echo",
							ChangePath: true,
						}},
					}, nil)
				t.argService.EXPECT().PrepareArgs(gomock.Any(), config.Operation{
					Name:       "before",
					Cmd:        "echo",
					ChangePath: true,
				}).Return([]string{"'Hello, World! before'"}, nil)
				t.operationService.EXPECT().GetOperationExecutionPath(gomock.Any(), "before").
					Return("", assert.AnError)
			},
			expectedErr: errors.New("failed to run before: failed to run before operation: before: failed to run command: failed to get operation execution path: assert.AnError general error for testing"),
		},
		"with error on prepare args": {
			preconditions: func(t *testController) {
				t.flagService.EXPECT().GetInitialFlags().
					Return(&entity.Flags{
						Operation: "operation",
					})
				t.operationService.EXPECT().GetEnhancedOperation(gomock.Any(), "operation").
					Return(config.Operation{
						Name: "operation",
						Cmd:  "echo",
					}, nil)
				t.argService.EXPECT().PrepareArgs(gomock.Any(), config.Operation{
					Name: "operation",
					Cmd:  "echo",
				}).Return([]string{}, assert.AnError)
			},
			expectedErr: errors.New("failed to prepare args: assert.AnError general error for testing"),
		},
		"with error on get enhanced operation": {
			preconditions: func(t *testController) {
				t.flagService.EXPECT().GetInitialFlags().
					Return(&entity.Flags{
						Operation: "operation",
					})
				t.operationService.EXPECT().GetEnhancedOperation(gomock.Any(), "operation").
					Return(config.Operation{}, assert.AnError)
			},
			expectedErr: errors.New("failed to get enhanced operation: assert.AnError general error for testing"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tc := newTestController(ctrl)
			testCase.preconditions(tc)

			service := tc.Build()
			err := service.Run(context.Background())

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})

	}
}

type testController struct {
	argService       *mocks.MockArgService
	operationService *mocks.MockOperationService
	flagService      *mocks.MockFlagService
}

func newTestController(ctrl *gomock.Controller) *testController {
	return &testController{
		flagService:      mocks.NewMockFlagService(ctrl),
		operationService: mocks.NewMockOperationService(ctrl),
		argService:       mocks.NewMockArgService(ctrl),
	}
}

func (t *testController) Build() *Service {
	return NewService(t.operationService,
		t.flagService,
		t.argService,
	)
}
