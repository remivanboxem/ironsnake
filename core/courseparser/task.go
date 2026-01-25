package courseparser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// EnvironmentLimits defines resource limits for Docker environments
type EnvironmentLimits struct {
	Time     string `yaml:"time"`      // Time limit in seconds
	HardTime string `yaml:"hard_time"` // Hard time limit (empty if not set)
	Memory   string `yaml:"memory"`    // Memory limit in MB
}

// EnvironmentParameters defines environment configuration
type EnvironmentParameters struct {
	Limits *EnvironmentLimits `yaml:"limits,omitempty"`
	RunCmd string             `yaml:"run_cmd,omitempty"`
}

// ProblemMap is a map of problem ID to Problem, with custom YAML unmarshaling
type ProblemMap map[string]Problem

// UnmarshalYAML handles polymorphic deserialization of problems
func (pm *ProblemMap) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("problems must be a mapping, got %v", node.Kind)
	}

	*pm = make(ProblemMap)

	// Iterate through key-value pairs
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		problemID := keyNode.Value

		// First, decode just the type field to determine the problem type
		var typeCheck struct {
			Type string `yaml:"type"`
		}
		if err := valueNode.Decode(&typeCheck); err != nil {
			return fmt.Errorf("problem %s: failed to read type: %w", problemID, err)
		}

		var problem Problem
		switch typeCheck.Type {
		case "code":
			var p CodeProblem
			if err := valueNode.Decode(&p); err != nil {
				return fmt.Errorf("problem %s (code): %w", problemID, err)
			}
			problem = &p

		case "multiple_choice":
			var p MultipleChoiceProblem
			if err := valueNode.Decode(&p); err != nil {
				return fmt.Errorf("problem %s (multiple_choice): %w", problemID, err)
			}
			problem = &p

		case "match":
			var p MatchProblem
			if err := valueNode.Decode(&p); err != nil {
				return fmt.Errorf("problem %s (match): %w", problemID, err)
			}
			problem = &p

		default:
			return fmt.Errorf("problem %s: unknown type %q", problemID, typeCheck.Type)
		}

		(*pm)[problemID] = problem
	}

	return nil
}

// TaskConfig represents a task.yaml file
type TaskConfig struct {
	Author                string                `yaml:"author"`
	ContactURL            string                `yaml:"contact_url"`
	Context               string                `yaml:"context"`
	EnvironmentID         string                `yaml:"environment_id"`
	EnvironmentType       string                `yaml:"environment_type"`
	EnvironmentParameters EnvironmentParameters `yaml:"environment_parameters"`
	File                  string                `yaml:"file"`
	Name                  string                `yaml:"name"`
	NetworkGrading        bool                  `yaml:"network_grading"`
	Problems              ProblemMap            `yaml:"problems"`

	// MCQ-specific fields
	InputRandom           int    `yaml:"input_random,omitempty"`
	RegenerateInputRandom string `yaml:"regenerate_input_random,omitempty"`
}

// IsDocker returns true if this task uses a Docker environment
func (t *TaskConfig) IsDocker() bool {
	return t.EnvironmentType == "docker"
}

// IsMCQ returns true if this task is a multiple choice quiz
func (t *TaskConfig) IsMCQ() bool {
	return t.EnvironmentType == "mcq"
}

// ParseTaskConfig parses a task.yaml file and returns a TaskConfig
func ParseTaskConfig(path string) (*TaskConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ParseError{
			File:    path,
			Message: "failed to read file",
			Err:     err,
		}
	}

	var config TaskConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, &ParseError{
			File:    path,
			Message: "failed to parse YAML",
			Err:     err,
		}
	}

	return &config, nil
}
