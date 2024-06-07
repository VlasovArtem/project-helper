package extractor

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"project-helper/internal/domain/entity"
)

func TestExtractTags(t *testing.T) {
	t.Parallel()

	service := NewService()

	tags := service.ExtractTags("${{tag1}} hello ${{tag2}}")

	assert.Len(t, tags, 2)
	assert.ElementsMatch(t, entity.Tags{"${{tag1}}", "${{tag2}}"}, tags)
}

func TestExtractTag(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		tag         entity.Tag
		expected    string
		expectedErr error
	}{
		"valid tag": {
			tag:      "${{tag1}}",
			expected: "tag1",
		},
		"invalid tag": {
			tag:         "tag1",
			expectedErr: errors.New("tag 'tag1' is not valid: tag value not found"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			service := NewService()

			tag, err := service.ExtractTag(testCase.tag)

			if testCase.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expected, tag)
			}
		})
	}

	service := NewService()

	tag, err := service.ExtractTag("${{tag1}}")

	assert.NoError(t, err)
	assert.Equal(t, "tag1", tag)
}
