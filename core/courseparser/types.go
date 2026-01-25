package courseparser

// ParsedCourse represents a fully loaded course from the filesystem
type ParsedCourse struct {
	DirPath  string                 // Source directory path
	CourseID string                 // Course ID (directory name, e.g., "CS01")
	Config   CourseConfig           // Parsed config.yaml
	Access   AccessConfig           // Parsed access.yaml
	Tasks    map[string]TaskConfig  // Task ID -> TaskConfig
	Syllabus *Syllabus              // Parsed syllabus (nil if not present)
}

// Problem is the interface for all problem types
type Problem interface {
	GetType() string
	GetName() string
	GetHeader() string
}

// BaseProblem contains fields common to all problem types
type BaseProblem struct {
	Type   string `yaml:"type"`
	Name   string `yaml:"name"`
	Header string `yaml:"header"`
}

func (p BaseProblem) GetType() string   { return p.Type }
func (p BaseProblem) GetName() string   { return p.Name }
func (p BaseProblem) GetHeader() string { return p.Header }

// CodeProblem represents a coding problem (type: "code")
type CodeProblem struct {
	BaseProblem `yaml:",inline"`
	Language    string `yaml:"language"`
	Default     string `yaml:"default"`
}

// Choice represents a single choice in a multiple choice question
type Choice struct {
	Text  string `yaml:"text"`
	Valid bool   `yaml:"valid"`
}

// MultipleChoiceProblem represents a multiple choice problem (type: "multiple_choice")
type MultipleChoiceProblem struct {
	BaseProblem `yaml:",inline"`
	Choices     []Choice `yaml:"choices"`
	Limit       int      `yaml:"limit"`
}

// MatchProblem represents a match/fill-in problem (type: "match")
type MatchProblem struct {
	BaseProblem `yaml:",inline"`
	Answer      string `yaml:"answer"`
}
