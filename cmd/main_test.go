//go:build integration

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
	"project-helper/internal/config"
	"project-helper/internal/domain/entity"
)

func TestIntegration(t *testing.T) {
	dir := t.TempDir()

	application := config.Application{
		Name: "project-helper-test",
		Path: dir,
		DynamicFlags: config.DynamicFlags{
			{
				Name:        "dynamic-flag",
				ShortName:   "d",
				Description: "dynamic flag description",
				Type:        entity.String,
			},
		},
		PredefinedArgs: config.PredefinedArgs{
			{
				Name: "predefined-tag-value",
				Type: entity.String,
				Args: config.Args{
					{
						Name:   "predefined-tag-name",
						Values: []string{"${{dynamic-flag}}.txt"},
					},
				},
			},
			{
				Name: "dynamic-flag",
				Type: entity.Array,
				Args: config.Args{
					{
						Name:   "predefined-tag-name",
						Values: []string{"test"},
					},
				},
			},
		},
		Operations: []config.Operation{
			{
				Name:        "test-echo",
				ShortName:   "te",
				Cmd:         "touch",
				Description: "test echo",
				ChangePath:  true,
				RunBefore: config.Operations{
					{
						Name: "test-before",
					},
				},
				PredefinedArgsTag: &config.PredefinedArgsTag{
					Name:  "dynamic-flag",
					Value: "predefined-tag-value",
				},
			},
			{
				Name:        "test-before",
				ShortName:   "tb",
				Cmd:         "touch",
				Description: "test before",
				ChangePath:  true,
				Args:        []string{"test-before.txt"},
			},
		},
	}

	applicationConfigFile, err := os.Create(dir + "/application.yaml")
	require.NoError(t, err)

	err = yaml.NewEncoder(applicationConfigFile).Encode(application)
	require.NoError(t, err)

	err = os.Setenv("CONFIG_PATH", applicationConfigFile.Name())
	require.NoError(t, err)

	os.Args = append(os.Args, "-o=test-echo", "--dynamic-flag=predefined-tag-name")

	main()

	_, err = os.ReadFile(dir + "/test-before.txt")
	require.NoError(t, err)

	_, err = os.ReadFile(dir + "/test.txt")
	require.NoError(t, err)
}
