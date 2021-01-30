package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type LogEvent struct {
	Date     time.Time `json:"date"`
	Duration int       `json:"duration"`
	Tag      string    `json:"tag"`
	Text     string    `json:"text,omitempty"`
}

// parseLog reads a log formated string and returns a LogEvent
func parseLog(line string) (LogEvent, error) {
	// Parse content
	// Log syntax is:
	//  YYYY-MM-DD   <dur>H  <tag>  <text>
	re := regexp.MustCompile(`^(\d{4}-[0-1]\d-[0-3]\d) +(\d+)H +([.a-z]+)(?:$| +(.*))`)
	content := re.FindStringSubmatch(line)
	if len(content) != 5 {
		return LogEvent{}, fmt.Errorf("Can't parse line with regexp %s", line)
	}

	// Read parsed values
	date, err := time.Parse("2006-01-02", content[1])
	if err != nil {
		return LogEvent{}, fmt.Errorf("Can't read date of the log %s: %w", content[1], err)
	}
	duration, err := strconv.Atoi(content[2])
	if err != nil {
		return LogEvent{}, fmt.Errorf("Can't read duration %s: %w", content[2], err)
	}

	// Create the log
	log := LogEvent{
		Date:     date,
		Duration: duration,
		Tag:      content[3],
		Text:     content[4],
	}

	return log, nil
}

// readLogFile reads fpath and creates a LogSet parsing every line
func readLogFile(fpath string) ([]LogEvent, error) {
	// Read file
	fmt.Errorf("Opening file %s\n", fpath)
	f, err := os.Open(fpath)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("Can't read file %s: %w", fpath, err)
	}
	scanner := bufio.NewScanner(f)

	// Parse each line of the file
	var logset []LogEvent
	for scanner.Scan() {
		log, err := parseLog(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("Can't read line %s: %w", fpath, err)
		}
		logset = append(logset, log)
	}

	return logset, nil
}

func main() {
	// Read data folder
	dirname := "data"

	d, err := os.Open(dirname)
	defer d.Close()
	if err != nil {
		println("Can't open folder %s: %w\n", dirname, err)
	}

	files, err := d.Readdir(-1)
	if err != nil {
		println("Can't read files in folder %s: %w\n", dirname, err)
	}

	// Parse each file as []LogEvent
	var alllogs []LogEvent
	for _, file := range files {
		filename := dirname + "/" + file.Name()
		data, err := readLogFile(filename)
		if err != nil {
			println("Can't parse file %s: %w\n", filename, err)
		}
		for _, d := range data {
			alllogs = append(alllogs, d)
		}
	}

	// Marshal as json
	jsondata, err := json.Marshal(alllogs)
	if err != nil {
		println("Can't output %s into JSON: %w", alllogs, err)
	}

	// Output
	os.Stdout.Write(jsondata)
}
