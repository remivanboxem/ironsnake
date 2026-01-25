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

// OrderedProblem wraps a Problem with its ID to maintain order
type OrderedProblem struct {
	ID      string
	Problem Problem
}

// ProblemMap is an ordered list of problems that preserves YAML order
type ProblemMap struct {
	Problems []OrderedProblem
	byID     map[string]Problem
}

// Get returns a problem by ID
func (pm *ProblemMap) Get(id string) (Problem, bool) {
	if pm.byID == nil {
		return nil, false
	}
	p, ok := pm.byID[id]
	return p, ok
}

// Len returns the number of problems
func (pm *ProblemMap) Len() int {
	return len(pm.Problems)
}

// UnmarshalYAML handles polymorphic deserialization of problems while preserving order
func (pm *ProblemMap) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("problems must be a mapping, got %v", node.Kind)
	}

	pm.Problems = make([]OrderedProblem, 0)
	pm.byID = make(map[string]Problem)

	// Iterate through key-value pairs (order is preserved in yaml.Node)
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

		pm.Problems = append(pm.Problems, OrderedProblem{ID: problemID, Problem: problem})
		pm.byID[problemID] = problem
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
