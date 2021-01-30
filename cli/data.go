package main

import (
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

// Tags returns a pointer on a map of tag names and the matching Logs in a
// LogSet.
func (logset *LogSet) Tags() *map[string]LogSet {
	tags := make(map[string]LogSet)
	for _, v := range *logset {
		tags[v.Tag] = append(tags[v.Tag], v)
	}
	return &tags
}
