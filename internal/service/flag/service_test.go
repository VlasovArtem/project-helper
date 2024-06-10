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
			DynamicFlags: map[string]*entity.DynamicFlagValue{
				"flag": {
					Value: utils.MakePointer("value"),
					Type:  entity.String,
					Name:  "flag",
				},
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
				DynamicFlags: map[string]*entity.DynamicFlagValue{
					"flag": {
						Value: utils.MakePointer("value"),
						Type:  entity.String,
						Name:  "flag",
					},
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
				DynamicFlags: map[string]*entity.DynamicFlagValue{
					"flag": {
						Value: utils.MakePointer("predefined_value"),
						Type:  entity.String,
						Name:  "flag",
					},
					"another_flag": {
						Value: utils.MakePointer("another_predefined_value"),
						Type:  entity.String,
						Name:  "another_flag",
					},
				},
			},
		},
		"with operation flags exists": {
			flagsService: &Service{
				initialFlags: &entity.Flags{
					DynamicFlags: map[string]*entity.DynamicFlagValue{
						"flag": {
							Value: utils.MakePointer("value"),
							Type:  entity.String,
							Name:  "flag",
						},
					},
				},
				operationFlags: map[string]*entity.Flags{
					"operation": {
						DynamicFlags: map[string]*entity.DynamicFlagValue{
							"flag": {
								Value: utils.MakePointer("operation_value"),
								Type:  entity.String,
								Name:  "flag",
							},
						},
					},
				},
			},
			operation: config.Operation{
				Name: "operation",
			},
			expectedFlags: &entity.Flags{
				DynamicFlags: map[string]*entity.DynamicFlagValue{
					"flag": {
						Value: utils.MakePointer("operation_value"),
						Type:  entity.String,
						Name:  "flag",
					},
				},
			},
		},
		"no predefined flags": {
			flagsService: NewFlagsService(&entity.Flags{
				DynamicFlags: map[string]*entity.DynamicFlagValue{
					"flag": {
						Value: utils.MakePointer("value"),
						Type:  entity.String,
						Name:  "flag",
					},
				},
			}),
			operation: config.Operation{},
			expectedFlags: &entity.Flags{
				DynamicFlags: map[string]*entity.DynamicFlagValue{
					"flag": {
						Value: utils.MakePointer("value"),
						Type:  entity.String,
						Name:  "flag",
					},
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
