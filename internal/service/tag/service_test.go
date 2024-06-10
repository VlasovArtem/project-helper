package tag

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"project-helper/internal/config"
	"project-helper/internal/domain/dto"
	"project-helper/internal/domain/entity"
	"project-helper/internal/service/tag/mocks"
	"project-helper/internal/utils"
)

func TestGetTagValue(t *testing.T) {
	t.Parallel()

	var (
		operation = config.Operation{
			Name: "operation",
		}
	)

	tests := map[string]struct {
		preconditions func(t *testController)
		input         *dto.GetTagValueRequest
		output        string
		expectedErr   error
	}{
		"success with string tag": {
			input: &dto.GetTagValueRequest{
				Flags: &entity.Flags{
					DynamicFlags: map[string]*entity.DynamicFlagValue{
						"tag1": {
							Name:  "tag1",
							Type:  entity.String,
							Value: utils.MakePointer("tag1—value"),
						},
					},
				},
				Operation:    operation,
				ExtractedTag: "tag1",
			},
			output: "tag1—value",
		},
		"success with string array tag": {
			input: &dto.GetTagValueRequest{
				Flags: &entity.Flags{
					DynamicFlags: map[string]*entity.DynamicFlagValue{
						"tag1": {
							Name:  "tag1",
							Type:  entity.Array,
							Value: &[]string{"tag1—value", "tag2—value"},
						},
					},
				},
				Operation:    operation,
				ExtractedTag: "tag1",
			},
			output: "tag1—value,tag2—value",
		},
		"success without pattern tag matches": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetAdditionalArgs().Return(map[string]string{
					"tag1": "tag1—value",
				})
			},
			input: &dto.GetTagValueRequest{
				Flags:        &entity.Flags{},
				Operation:    operation,
				ExtractedTag: "tag1",
			},
			output: "tag1—value",
		},
		"success without pattern tag matches and execution-path tag": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetAdditionalArgs().Return(map[string]string{
					"tag1": "tag1—value",
				})
				t.configService.EXPECT().GetApplicationPath().Return("application-path")
			},
			input: &dto.GetTagValueRequest{
				Flags: &entity.Flags{},
				Operation: config.Operation{
					ExecutionPath: "execution-path",
					ChangePath:    true,
				},
				ExtractedTag: entity.ExecutionPathTag,
			},
			output: "application-path/execution-path",
		},
		"without pattern tag matches and no additional tag": {
			preconditions: func(t *testController) {
				t.configService.EXPECT().GetAdditionalArgs().Return(map[string]string{})
			},
			input: &dto.GetTagValueRequest{
				Flags:        &entity.Flags{},
				Operation:    operation,
				ExtractedTag: "tag1",
			},
			expectedErr: errors.New("failed to check additional args: additional arg tag1 not found: additional arg not found"),
		},
		"with invalid request": {
			preconditions: func(t *testController) {},
			input:         &dto.GetTagValueRequest{},
			expectedErr:   errors.New("request is not valid"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			controller := newTestController(ctrl)

			if testCase.preconditions != nil {
				testCase.preconditions(controller)
			}

			service := controller.Build()

			output, err := service.GetTagValue(testCase.input)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.output, output)
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
