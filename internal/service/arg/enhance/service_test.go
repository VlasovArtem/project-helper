package enhance

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"project-helper/internal/config"
	"project-helper/internal/domain/dto"
	"project-helper/internal/domain/entity"
	"project-helper/internal/service/arg/enhance/mocks"
)

func TestEnhanceArgs(t *testing.T) {
	t.Parallel()

	var (
		operation = config.Operation{
			Name: "operation",
		}
		flags = entity.Flags{}
	)

	tests := map[string]struct {
		precondition   func(*testController)
		input          *dto.EnhanceArgsRequest
		expectedOutput []string
		expectedError  error
	}{
		"success": {
			precondition: func(tc *testController) {
				tc.extractorService.EXPECT().ExtractTags(entity.Arg("arg1 - tag1")).
					Return(entity.Tags{"tag1"})
				tc.extractorService.EXPECT().ExtractTags(entity.Arg("arg2")).
					Return(entity.Tags{})
				tc.extractorService.EXPECT().ExtractTag(entity.Tag("tag1")).
					Return("tag_value1", nil)
				tc.tagService.EXPECT().GetTagValue(&dto.GetTagValueRequest{
					Operation:    operation,
					Flags:        &flags,
					ExtractedTag: "tag_value1",
				}).Return("new_tag_value1", nil)
				tc.predefinedArgService.EXPECT().TryToFindPredefinedArgValue(&dto.TryToFindPredefinedArgRequest{
					ParsedTag: "tag_value1",
					Value:     "new_tag_value1",
				}).Return(`"quoted_new_tag_value1"`, nil)
			},
			input: &dto.EnhanceArgsRequest{
				Operation: operation,
				Flags:     &flags,
				Args:      []string{"arg1 - tag1", "arg2"},
			},
			expectedOutput: []string{"arg1 - \"quoted_new_tag_value1\"", "arg2"},
		},
		"with error on try to find predefined arg value": {
			precondition: func(tc *testController) {
				tc.extractorService.EXPECT().ExtractTags(entity.Arg("arg1 - tag1")).
					Return(entity.Tags{"tag1"})
				tc.extractorService.EXPECT().ExtractTag(entity.Tag("tag1")).
					Return("tag_value1", nil)
				tc.tagService.EXPECT().GetTagValue(&dto.GetTagValueRequest{
					Operation:    operation,
					Flags:        &flags,
					ExtractedTag: "tag_value1",
				}).Return("new_tag_value1", nil)
				tc.predefinedArgService.EXPECT().TryToFindPredefinedArgValue(&dto.TryToFindPredefinedArgRequest{
					ParsedTag: "tag_value1",
					Value:     "new_tag_value1",
				}).Return("", assert.AnError)
			},
			input: &dto.EnhanceArgsRequest{
				Operation: operation,
				Flags:     &flags,
				Args:      []string{"arg1 - tag1", "arg2"},
			},
			expectedError: errors.New("failed to try to find predefined arg: assert.AnError general error for testing"),
		},
		"with error on get tag value": {
			precondition: func(tc *testController) {
				tc.extractorService.EXPECT().ExtractTags(entity.Arg("arg1 - tag1")).
					Return(entity.Tags{"tag1"})
				tc.extractorService.EXPECT().ExtractTag(entity.Tag("tag1")).
					Return("tag_value1", nil)
				tc.tagService.EXPECT().GetTagValue(&dto.GetTagValueRequest{
					Operation:    operation,
					Flags:        &flags,
					ExtractedTag: "tag_value1",
				}).Return("", assert.AnError)
			},
			input: &dto.EnhanceArgsRequest{
				Operation: operation,
				Flags:     &flags,
				Args:      []string{"arg1 - tag1", "arg2"},
			},
			expectedError: errors.New("failed to get tag value: assert.AnError general error for testing"),
		},
		"with error on extract tag": {
			precondition: func(tc *testController) {
				tc.extractorService.EXPECT().ExtractTags(entity.Arg("arg1 - tag1")).
					Return(entity.Tags{"tag1"})
				tc.extractorService.EXPECT().ExtractTag(entity.Tag("tag1")).
					Return("", assert.AnError)
			},
			input: &dto.EnhanceArgsRequest{
				Operation: operation,
				Flags:     &flags,
				Args:      []string{"arg1 - tag1", "arg2"},
			},
			expectedError: errors.New("failed to extract tag: assert.AnError general error for testing"),
		},
		"with invalid request": {
			input:         &dto.EnhanceArgsRequest{},
			expectedError: errors.New("request is not valid"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			controller := newTestController(ctrl)

			if testCase.precondition != nil {
				testCase.precondition(controller)
			}

			service := controller.Build()

			actual, err := service.EnhanceArgs(testCase.input)

			if testCase.expectedError != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedOutput, actual)
			}
		})
	}
}

func TestGetEnhancedOperationArgs(t *testing.T) {
	t.Parallel()

	var (
		operation = config.Operation{
			Name: "operation",
			Args: []string{"arg1", "arg2"},
			PredefinedArgsTag: &config.PredefinedArgsTag{
				Name: "predefined_arg",
			},
		}
		flags = entity.Flags{}
	)

	tests := map[string]struct {
		precondition   func(*testController)
		input          *dto.GetEnhancedOperationArgs
		expectedOutput []string
		expectedError  error
	}{
		"success": {
			precondition: func(tc *testController) {
				tc.extractorService.EXPECT().ExtractTags(entity.Arg("arg1")).
					Return(entity.Tags{"tag1"})
				tc.extractorService.EXPECT().ExtractTag(entity.Tag("tag1")).
					Return("predefined_arg", nil)
				tc.predefinedArgService.EXPECT().GetPredefinedArgValues(&dto.GetPredefinedArgsRequest{
					Flags:             &flags,
					PredefinedArgsTag: &config.PredefinedArgsTag{Name: "predefined_arg"},
				}).Return([]string{"predefined_arg1", "predefined_arg2"}, nil)

				tc.extractorService.EXPECT().ExtractTags(entity.Arg("arg2")).
					Return(entity.Tags{})
			},
			input: &dto.GetEnhancedOperationArgs{
				Operation: operation,
				Flags:     &flags,
			},
			expectedOutput: []string{"predefined_arg1", "predefined_arg2", "arg2"},
		},
		"with not predefined args": {
			precondition: func(tc *testController) {
				tc.extractorService.EXPECT().ExtractTags(entity.Arg("arg1")).
					Return(entity.Tags{"tag1"})
				tc.extractorService.EXPECT().ExtractTag(entity.Tag("tag1")).
					Return("predefined_arg", nil)
				tc.predefinedArgService.EXPECT().GetPredefinedArgValues(&dto.GetPredefinedArgsRequest{
					Flags:             &flags,
					PredefinedArgsTag: &config.PredefinedArgsTag{Name: "predefined_arg"},
				}).Return([]string{}, nil)

				tc.extractorService.EXPECT().ExtractTags(entity.Arg("arg2")).
					Return(entity.Tags{})
			},
			input: &dto.GetEnhancedOperationArgs{
				Operation: operation,
				Flags:     &flags,
			},
			expectedOutput: []string{"arg1", "arg2"},
		},
		"with get predefined args values error": {
			precondition: func(tc *testController) {
				tc.extractorService.EXPECT().ExtractTags(entity.Arg("arg1")).
					Return(entity.Tags{"tag1"})
				tc.extractorService.EXPECT().ExtractTag(entity.Tag("tag1")).
					Return("predefined_arg", nil)
				tc.predefinedArgService.EXPECT().GetPredefinedArgValues(&dto.GetPredefinedArgsRequest{
					Flags:             &flags,
					PredefinedArgsTag: &config.PredefinedArgsTag{Name: "predefined_arg"},
				}).Return([]string{}, assert.AnError)

			},
			input: &dto.GetEnhancedOperationArgs{
				Operation: operation,
				Flags:     &flags,
			},
			expectedError: errors.New("failed to get predefined args: assert.AnError general error for testing"),
		},
		"with extract tag error": {
			precondition: func(tc *testController) {
				tc.extractorService.EXPECT().ExtractTags(entity.Arg("arg1")).
					Return(entity.Tags{"tag1"})
				tc.extractorService.EXPECT().ExtractTag(entity.Tag("tag1")).
					Return("", assert.AnError)
			},
			input: &dto.GetEnhancedOperationArgs{
				Operation: operation,
				Flags:     &flags,
			},
			expectedError: errors.New("failed to extract tag: assert.AnError general error for testing"),
		},
		"with invalid request": {
			input:         &dto.GetEnhancedOperationArgs{},
			expectedError: errors.New("request is not valid"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			controller := newTestController(ctrl)

			if testCase.precondition != nil {
				testCase.precondition(controller)
			}

			service := controller.Build()

			actual, err := service.GetEnhancedOperationArgs(testCase.input)

			if testCase.expectedError != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedOutput, actual)
			}
		})
	}
}

type testController struct {
	extractorService     *mocks.MockExtractorService
	tagService           *mocks.MockTagService
	predefinedArgService *mocks.MockPredefinedArgService
}

func newTestController(ctrl *gomock.Controller) *testController {
	return &testController{
		extractorService:     mocks.NewMockExtractorService(ctrl),
		tagService:           mocks.NewMockTagService(ctrl),
		predefinedArgService: mocks.NewMockPredefinedArgService(ctrl),
	}
}

func (t *testController) Build() *Service {
	return NewService(
		t.extractorService,
		t.tagService,
		t.predefinedArgService,
	)
}
