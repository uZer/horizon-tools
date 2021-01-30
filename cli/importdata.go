package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

// parseLog reads a log formated string and returns a Log
func parseLog(line string) (Log, error) {
	// Parse content
	// Log syntax is:
	//  YYYY-MM-DD   <dur>H  <tag>  <text>
	re := regexp.MustCompile(`^(\d{4}-[0-1]\d-[0-3]\d) +(\d+)H +([.a-z]+)(?:$| +(.*))`)
	content := re.FindStringSubmatch(line)
	if len(content) != 5 {
		return Log{}, fmt.Errorf("can't parse line with regexp %s", line)
	}

	// Read parsed values
	date, err := time.Parse("2006-01-02", content[1])
	if err != nil {
		return Log{}, fmt.Errorf("can't read date of the log %s: %w", content[1], err)
	}
	duration, err := strconv.Atoi(content[2])
	if err != nil {
		return Log{}, fmt.Errorf("can't read duration %s: %w", content[2], err)
	}

	// Create the log
	log := Log{
		Date:     date,
		Duration: duration,
		Tag:      content[3],
		Text:     content[4],
	}

	return log, nil
}

// readLogFile reads fpath and creates a LogSet parsing every line
func readLogFile(fpath string) (LogSet, error) {
	// Read file
	fmt.Errorf("opening file %s\n", fpath)
	f, err := os.Open(fpath)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("can't read file %s: %w", fpath, err)
	}
	scanner := bufio.NewScanner(f)

	// Parse each line of the file
	var logset LogSet
	for scanner.Scan() {
		log, err := parseLog(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("can't read line %s: %w", fpath, err)
		}
		logset = append(logset, log)
	}

	return logset, nil
}
