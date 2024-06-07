# Project Helper

The project helper allow to create an application that allow to run preconfigured commands to help in the development of
a project.

## Overview

Project Helper is a Go-based utility designed to streamline and automate various project-related tasks. It leverages
configuration files, command-line flags, and predefined arguments to execute operations efficiently.

## Getting Started

### Prerequisites

* Go 1.22
* make (optional, for running tests)

### Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/project-helper.git
cd project-helper
```

Install dependencies

```bash
go mod tidy
```

### Configuration

The application expects a configuration file in YAML format. The default path for the configuration file
is `$XDG_CONFIG_HOME/project-helper/application.yaml`. You can override this path by setting the `CONFIG_PATH`
environment variable.
Example configuration (`application.yaml`):

```yaml
name: "Project Helper"
path: "/path/to/project"
dynamicFlags:
  - name: "flag"
    shortName: "f"
    description: "A dynamic flag"
    type: "string"
    default: "default_value"
patternTags:
  - name: "tag1"
    type: "string"
    join: ","
predefinedArgs:
  - name: "arg1"
    type: "string"
    args:
      - name: "value1"
        values: [ "val1", "val2" ]
operations:
  - name: "operation1"
    shortName: "op1"
    description: "An example operation"
    cmd: "echo"
    args: [ "Hello, World!" ]
    executionPath: "/path/to/execute"
    changePath: true
    predefinedArgsTag:
      name: "arg1"
      value: "value1"
    runBefore:
      - name: "setup"
        cmd: "echo"
        args: [ "Setting up..." ]
```

### Running the Application

To run the application, use the following command:

```bash
go run cmd/main.go --operation=operation1
```

### Unit Tests

To run the unit tests, you can use the `go test` command or `make` if you have a Makefile set up.

```shell
go test ./...
```

or

```shell
make unit-test
```

## Project Structure

`cmd/main.go`
The entry point of the application. It initializes services and runs the main operation.
`internal/config`
Contains configuration-related code, including the structure of the configuration file and methods to parse it.
`internal/domain`
Defines the core domain entities, DTOs (Data Transfer Objects), and custom errors used throughout the application.
`internal/service`
Contains the business logic of the application, organized into various services such
as `arg`, `config`, `flag`, `operation`, projecthelper, and tag.
`internal/utils`
Utility functions that are used across different parts of the application.

## Key Components

### Services

* `Arg Service`: Handles argument preparation and enhancement.
* `Config Service`: Manages application configuration.
* `Flag Service`: Parses and validates command-line flags.
* `Operation Service`: Retrieves and enhances operations.
* `Project Helper Service`: Orchestrates the execution of operations.
* `Tag Service`: Extracts and processes tags from arguments.

### Mocks

Mocks are generated using `MockGen` for unit testing purposes. They are located in the `mocks` subdirectories within
each service directory.

## Example Usage

1. Define your operations and configurations in `application.yaml`.
2. Run the application with the desired operation.
3. The application will parse the configuration, enhance arguments, and execute the specified commands.

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.