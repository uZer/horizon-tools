package main

import (
	"encoding/json"
	"os"
)

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
	alllogs := make(LogSet)
	logId := 0
	for _, file := range files {
		filename := dirname + "/" + file.Name()
		data, err := readLogFile(filename)
		if err != nil {
			println("Can't parse file %s: %w\n", filename, err)
		}
		for _, d := range data {
			alllogs[logId] = d
			logId++
		}
	}

	// Marshal as json
	jsonLogs, err := json.Marshal(alllogs)
	if err != nil {
		println("Error marshaling LogSet in JSON: %w", err)
	}
	jsonTags, err := json.Marshal(alllogs.Tags())
	if err != nil {
		println("Error marshaling LogSet Tags in JSON: %w", err)
	}

	// Output
	generateOutputs()

	os.Stdout.Write(jsonLogs)
	os.Stdout.Write(jsonTags)
}
