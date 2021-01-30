package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Log is a single parsed log entry
type Log struct {
	Date     time.Time `json:"date"`
	Duration int       `json:"duration"`
	Tag      string    `json:"tag"`
	Text     string    `json:"text,omitempty"`
}

// LogSet is a slice of Logs
type LogSet []Log

// IndexLog reads a formated string and appends it in the LogSet as a Log
func (l *LogSet) ImportLog(line string) error {
	// Parse content
	// Log syntax is:
	//  YYYY-MM-DD   <dur>H  <tag>  <text>
	re := regexp.MustCompile(`^(\d{4}-[0-1]\d-[0-3]\d) +(\d+)H +([.a-z]+)(?:$| +(.*))`)
	content := re.FindStringSubmatch(line)
	if len(content) != 5 {
		return fmt.Errorf("can't parse line with regexp %s", line)
	}

	// Read parsed values
	date, err := time.Parse("2006-01-02", content[1])
	if err != nil {
		return fmt.Errorf("can't read date of the log %s: %w", content[1], err)
	}
	duration, err := strconv.Atoi(content[2])
	if err != nil {
		return fmt.Errorf("can't read duration %s: %w", content[2], err)
	}

	// Create the log
	log := Log{
		Date:     date,
		Duration: duration,
		Tag:      content[3],
		Text:     content[4],
	}

	// Append value to self
	*l = append(*l, log)
	return nil
}

// ImportFile reads fpath and insert every line as log in the LogSet
func (l *LogSet) ImportFile(fpath string) error {
	// Read file
	fmt.Errorf("opening file %s", fpath)
	f, err := os.Open(fpath)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("can't read file %s: %w", fpath, err)
	}
	scanner := bufio.NewScanner(f)

	// Parse each line of the file
	for scanner.Scan() {
		err := l.ImportLog(scanner.Text())
		if err != nil {
			return fmt.Errorf("can't insert line %s: %w", scanner.Text(), err)
		}
	}

	return nil
}

// Tags returns a slice of tag names contained in LogSet
func (l *LogSet) Tags() []string {

	// Build an intermediary array indexing tags as keys
	uniqueTags := make(map[string]struct{})
	for _, v := range *l {
		uniqueTags[v.Tag] = struct{}{}
	}

	// Return a tag slice
	tags := []string{}
	for k := range uniqueTags {
		tags = append(tags, k)
	}

	return tags
}
