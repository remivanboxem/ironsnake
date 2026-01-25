package courseparser

import (
	"os"

	"gopkg.in/yaml.v3"
)

// CourseConfig represents the config.yaml file for a course
type CourseConfig struct {
	Accessible           bool           `yaml:"accessible"`
	Name                 string         `yaml:"name"`
	Admins               []string       `yaml:"admins"`
	Tutors               []string       `yaml:"tutors"`
	GroupsStudentChoice  bool           `yaml:"groups_student_choice"`
	AllowUnregister      bool           `yaml:"allow_unregister"`
	AllowPreview         bool           `yaml:"allow_preview"`
	Registration         bool           `yaml:"registration"`
	RegistrationPassword *string        `yaml:"registration_password"`
	RegistrationAC       *string        `yaml:"registration_ac"`
	RegistrationACAccept bool           `yaml:"registration_ac_accept"`
	RegistrationACList   []string       `yaml:"registration_ac_list"`
	Tags                 map[string]any `yaml:"tags"`
}

// ParseCourseConfig parses a config.yaml file and returns a CourseConfig
func ParseCourseConfig(path string) (*CourseConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &ParseError{
			File:    path,
			Message: "failed to read file",
			Err:     err,
		}
	}

	var config CourseConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, &ParseError{
			File:    path,
			Message: "failed to parse YAML",
			Err:     err,
		}
	}

	return &config, nil
}
