package courseparser

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const dateTimeLayout = "2006-01-02 15:04:05"

// TaskAccessibility represents the accessibility field which can be a bool or a date range
type TaskAccessibility struct {
	IsBoolean bool                      // True if this is a boolean value
	BoolValue bool                      // The boolean value (if IsBoolean is true)
	DateRange *AccessibilityDateRange   // The date range (if IsBoolean is false)
}

// AccessibilityDateRange represents a date range for task accessibility
// Format: "start/deadline/soft_deadline"
type AccessibilityDateRange struct {
	Start        time.Time // When the task becomes accessible
	Deadline     time.Time // Hard deadline
	SoftDeadline time.Time // Soft deadline (typically before hard deadline)
}

// UnmarshalYAML handles the polymorphic accessibility field
func (a *TaskAccessibility) UnmarshalYAML(node *yaml.Node) error {
	// Try boolean first
	var boolVal bool
	if err := node.Decode(&boolVal); err == nil {
		a.IsBoolean = true
		a.BoolValue = boolVal
		return nil
	}

	// Try date range string
	var strVal string
	if err := node.Decode(&strVal); err != nil {
		return fmt.Errorf("accessibility must be bool or date range string, got: %v", node.Value)
	}

	// Parse date range: "2026-01-25 19:15:03/2026-01-29 19:15:07/2026-01-28 19:15:04"
	parts := strings.Split(strVal, "/")
	if len(parts) != 3 {
		return fmt.Errorf("date range must have 3 parts (start/deadline/soft_deadline), got %d: %s", len(parts), strVal)
	}

	start, err := time.Parse(dateTimeLayout, parts[0])
	if err != nil {
		return fmt.Errorf("invalid start date %q: %w", parts[0], err)
	}

	deadline, err := time.Parse(dateTimeLayout, parts[1])
	if err != nil {
		return fmt.Errorf("invalid deadline %q: %w", parts[1], err)
	}

	softDeadline, err := time.Parse(dateTimeLayout, parts[2])
	if err != nil {
		return fmt.Errorf("invalid soft deadline %q: %w", parts[2], err)
	}

	a.IsBoolean = false
	a.DateRange = &AccessibilityDateRange{
		Start:        start,
		Deadline:     deadline,
		SoftDeadline: softDeadline,
	}

	return nil
}

// IsAccessible returns whether the task is currently accessible
func (a *TaskAccessibility) IsAccessible() bool {
	if a.IsBoolean {
		return a.BoolValue
	}
	if a.DateRange == nil {
		return false
	}
	now := time.Now()
	return now.After(a.DateRange.Start) && now.Before(a.DateRange.Deadline)
}

// SubmissionLimit defines rate limiting for task submissions
type SubmissionLimit struct {
	Amount int `yaml:"amount"` // Maximum submissions allowed
	Period int `yaml:"period"` // Time period in minutes
}

// TaskAccessConfig represents the configuration for a single task's access
type TaskAccessConfig struct {
	Accessibility       TaskAccessibility `yaml:"accessibility"`
	EvaluationMode      string            `yaml:"evaluation_mode,omitempty"`      // "best" or "last"
	NoStoredSubmissions int               `yaml:"no_stored_submissions,omitempty"` // Max stored submissions
	SubmissionLimit     *SubmissionLimit  `yaml:"submission_limit,omitempty"`
}

// DispenserData represents the dispenser_data section
type DispenserData struct {
	Config    map[string]TaskAccessConfig `yaml:"config"`
	Imported  bool                        `yaml:"imported"`
	Converted bool                        `yaml:"converted"`
}

// AccessConfig represents the access.yaml file
type AccessConfig struct {
	DispenserData DispenserData `yaml:"dispenser_data"`
}

// ParseAccessConfig parses an access.yaml file and returns an AccessConfig
func ParseAccessConfig(path string) (*AccessConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ParseError{
			File:    path,
			Message: "failed to read file",
			Err:     err,
		}
	}

	var config AccessConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, &ParseError{
			File:    path,
			Message: "failed to parse YAML",
			Err:     err,
		}
	}

	return &config, nil
}
