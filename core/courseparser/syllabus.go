package courseparser

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

// BookConfig represents the book.toml file
type BookConfig struct {
	Book BookInfo `toml:"book"`
}

// BookInfo contains the book metadata
type BookInfo struct {
	Title  string `toml:"title"`
	Author string `toml:"author"`
}

// SummaryEntry represents a single entry in the SUMMARY.md table of contents
type SummaryEntry struct {
	Title    string         // Display title
	Path     string         // Relative path to the markdown file
	Level    int            // Nesting level (0 = chapter heading, 1+ = content)
	Children []SummaryEntry // Nested entries
}

// Syllabus represents the complete syllabus structure
type Syllabus struct {
	Book    BookConfig     // Parsed book.toml
	Summary []SummaryEntry // Parsed SUMMARY.md entries
}

// ParseBookConfig parses a book.toml file
func ParseBookConfig(path string) (*BookConfig, error) {
	var config BookConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, &ParseError{
			File:    path,
			Message: "failed to parse TOML",
			Err:     err,
		}
	}
	return &config, nil
}

// linkPattern matches markdown links: [Title](path.md)
var linkPattern = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

// chapterPattern matches chapter headings: # Chapter Name
var chapterPattern = regexp.MustCompile(`^#\s+(.+)$`)

// ParseSummary parses a SUMMARY.md file and returns the table of contents entries
func ParseSummary(path string) ([]SummaryEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, &ParseError{
			File:    path,
			Message: "failed to open file",
			Err:     err,
		}
	}
	defer file.Close()

	var entries []SummaryEntry
	var currentChapter *SummaryEntry

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and the "Table of Contents" header
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "# Table of Contents") {
			continue
		}

		// Check for chapter heading
		if matches := chapterPattern.FindStringSubmatch(line); matches != nil {
			// Save previous chapter if exists
			if currentChapter != nil {
				entries = append(entries, *currentChapter)
			}
			currentChapter = &SummaryEntry{
				Title:    matches[1],
				Level:    0,
				Children: []SummaryEntry{},
			}
			continue
		}

		// Check for markdown link
		if matches := linkPattern.FindStringSubmatch(line); matches != nil {
			entry := SummaryEntry{
				Title: matches[1],
				Path:  matches[2],
				Level: countLeadingSpaces(line)/2 + 1, // Calculate nesting level
			}

			if currentChapter != nil {
				currentChapter.Children = append(currentChapter.Children, entry)
			} else {
				entries = append(entries, entry)
			}
		}
	}

	// Don't forget the last chapter
	if currentChapter != nil {
		entries = append(entries, *currentChapter)
	}

	if err := scanner.Err(); err != nil {
		return nil, &ParseError{
			File:    path,
			Message: "failed to read file",
			Err:     err,
		}
	}

	return entries, nil
}

// countLeadingSpaces counts the number of leading spaces (or tabs as 2 spaces)
func countLeadingSpaces(s string) int {
	count := 0
	for _, ch := range s {
		if ch == ' ' {
			count++
		} else if ch == '\t' {
			count += 2
		} else if ch == '-' {
			// List item marker, stop counting
			break
		} else {
			break
		}
	}
	return count
}

// ParseSyllabus parses both book.toml and SUMMARY.md from a syllabus directory
func ParseSyllabus(syllabusDir string) (*Syllabus, error) {
	bookPath := syllabusDir + "/book.toml"
	summaryPath := syllabusDir + "/SUMMARY.md"

	bookConfig, err := ParseBookConfig(bookPath)
	if err != nil {
		return nil, err
	}

	summary, err := ParseSummary(summaryPath)
	if err != nil {
		return nil, err
	}

	return &Syllabus{
		Book:    *bookConfig,
		Summary: summary,
	}, nil
}
