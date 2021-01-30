package main

import (
	"time"
)

type LogEvent struct {
	Date     time.Time `json:"date"`
	Duration int       `json:"duration"`
	Tag      string    `json:"tag"`
	Text     string    `json:"text,omitempty"`
}

type LogSet map[int]LogEvent

func (logset *LogSet) Tags() *map[string][]int {
	tags := make(map[string][]int)
	for k, v := range *logset {
		tags[v.Tag] = append(tags[v.Tag], k)
	}
	return &tags
}
