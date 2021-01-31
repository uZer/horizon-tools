package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

// ImportLog reads a formated string and appends it in the LogSet as a Log
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
	println("opening file %s", fpath)
	f, err := os.Open(fpath)
	if err != nil {
		return fmt.Errorf("can't read file %s: %w", fpath, err)
	}
	defer f.Close()
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
	for _, log := range *l {
		uniqueTags[log.Tag] = struct{}{}
	}

	// Return a tag slice
	tags := []string{}
	for k := range uniqueTags {
		tags = append(tags, k)
	}

	return tags
}

// Filter is a method to extract a subset of entries from the current LogSet.
// Each log is tested against the function passed in parameters. If the result
// of the function is true, the log is added to the output
// TODO: use interface
func (l *LogSet) Filter(test func(*Log, []string) bool, args []string) LogSet {
	set := LogSet{}
	for _, log := range *l {
		if test(&log, args) {
			set = append(set, log)
		}
	}
	return set
}

// HasTag returns true if the string passed as an arguments is equal to the
// Log's tag.
func HasTag(l *Log, tag []string) bool { return l.Tag == tag[0] }

// HasDay returns true if the Log's date matches the date passed as argument.
// The date format is YYYY-MM-DD
func HasDay(l *Log, date []string) bool {
	day, _ := time.Parse("2006-01-02", date[0])
	return l.Date.Equal(day)
}

// HasWeek returns true if the Log's week matches the date passed as argument.
// The week format is using two strings: year and week number between 1 and 53
func HasWeek(l *Log, date []string) bool {
	year, _ := strconv.Atoi(date[0])
	week, _ := strconv.Atoi(date[1])
	dy, dw := l.Date.ISOWeek()
	return dy == year && dw == week
}

// HasMonth returns true if the Log's date match the date passed as argument.
// The date format is YYYY-MM
// FIXME: doesn't seem to work
func HasMonth(l *Log, date []string) bool {
	time, _ := time.Parse("2006-01", date[0])
	return l.Date.Month() == time.Month()
}

// HasYear returns true if the Log's date match the year passed as argument.
// The date format is YYYY
func HasYear(l *Log, date []string) bool {
	year, _ := strconv.Atoi(date[0])
	return l.Date.Year() == year
}

// Contains returns true if the text of a log contains the substring passed as
// an argument. This method is not case sensitive.
func Contains(l *Log, text []string) bool {
	return strings.Contains(l.Text, strings.ToLower(text[0]))
}
