package config

import (
	"github.com/pkg/errors"
	"project-helper/internal/domain/entity"
)

type Application struct {
	Name           string
	Operations     Operations
	Path           string
	DynamicFlags   DynamicFlags   `yaml:"dynamicFlags"`
	PredefinedArgs PredefinedArgs `yaml:"predefinedArgs"`
}

type Operations []Operation

type Operation struct {
	Description       string
	Name              string
	ShortName         string `yaml:"shortName"`
	Cmd               string
	Args              []string
	ExecutionPath     string             `yaml:"executionPath"`
	ChangePath        bool               `yaml:"changePath"`
	PredefinedArgsTag *PredefinedArgsTag `yaml:"predefinedArgsTag"`
	RunBefore         Operations         `yaml:"runBefore"`
	PredefinedFlags   PredefinedFlags    `yaml:"predefinedFlags"`
}

type PredefinedArgsTag struct {
	Name  string
	Value string
}

type PredefinedFlags []PredefinedFlag

type PredefinedFlag struct {
	Name  string
	Value string
}

func (a *Application) GetOperationsMap() map[string]Operation {
	operationsMap := make(map[string]Operation)
	for _, operation := range a.Operations {
		if operation.Name != "" {
			operationsMap[operation.Name] = operation
		}

		if operation.ShortName != "" {
			operationsMap[operation.ShortName] = operation
		}
	}
	return operationsMap
}

func (a *Application) GetPredefinedArgs() map[string]PredefinedArg {
	predefinedArgs := make(map[string]PredefinedArg)
	for _, predefinedArg := range a.PredefinedArgs {
		if predefinedArg.Name != "" {
			predefinedArgs[predefinedArg.Name] = predefinedArg
		}
	}
	return predefinedArgs
}

type DynamicFlags []DynamicFlag

type DynamicFlag struct {
	Name        string
	ShortName   string `yaml:"shortName"`
	Description string
	Type        entity.Type
	Default     string
}

type PredefinedArgs []PredefinedArg

type PredefinedArg struct {
	Name string
	Type entity.Type
	Args Args
}

type Args []Arg

func (a Args) GetArgValues(name string) ([]string, error) {
	for _, arg := range a {
		if arg.Name == name {
			return arg.Values, nil
		}
	}
	return nil, errors.Errorf("arg %s not found", name)
}

type Arg struct {
	Name   string
	Values []string
}
