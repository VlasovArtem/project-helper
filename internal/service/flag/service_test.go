package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"project-helper/internal/config"
	"project-helper/internal/domain/entity"
	"project-helper/internal/utils"
)

func TestFlagsServiceGetInitialFlags(t *testing.T) {
	t.Parallel()

	service := Service{
		initialFlags: &entity.Flags{
			DynamicFlags: map[string]any{
				"flag": utils.MakePointer("value"),
			},
		},
	}

	flags := service.GetInitialFlags()

	assert.Equal(t, "value", flags.GetFlagStringValue("flag"))
}

func TestFlagsServiceGetOperationFlags(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		flagsService  *Service
		operation     config.Operation
		expectedFlags *entity.Flags
	}{
		"no operation flags exists": {
			flagsService: NewFlagsService(&entity.Flags{
				DynamicFlags: map[string]any{
					"flag": utils.MakePointer("value"),
				},
			}),
			operation: config.Operation{
				Name: "operation",
				PredefinedFlags: config.PredefinedFlags{
					{
						Name:  "flag",
						Value: "predefined_value",
					},
					{
						Name:  "another_flag",
						Value: "another_predefined_value",
					},
				},
			},
			expectedFlags: &entity.Flags{
				DynamicFlags: map[string]any{
					"flag":         utils.MakePointer("predefined_value"),
					"another_flag": utils.MakePointer("another_predefined_value"),
				},
			},
		},
		"with operation flags exists": {
			flagsService: &Service{
				initialFlags: &entity.Flags{
					DynamicFlags: map[string]any{
						"flag": utils.MakePointer("value"),
					},
				},
				operationFlags: map[string]*entity.Flags{
					"operation": {
						DynamicFlags: map[string]any{
							"flag": utils.MakePointer("operation_value"),
						},
					},
				},
			},
			operation: config.Operation{
				Name: "operation",
			},
			expectedFlags: &entity.Flags{
				DynamicFlags: map[string]any{
					"flag": utils.MakePointer("operation_value"),
				},
			},
		},
		"no predefined flags": {
			flagsService: NewFlagsService(&entity.Flags{
				DynamicFlags: map[string]any{
					"flag": utils.MakePointer("value"),
				},
			}),
			operation: config.Operation{},
			expectedFlags: &entity.Flags{
				DynamicFlags: map[string]any{
					"flag": utils.MakePointer("value"),
				},
			},
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			flags := testCase.flagsService.GetOperationFlags(testCase.operation)

			assert.Equal(t, testCase.expectedFlags, flags)
		})
	}
}
