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
		println("can't open folder %s: %w\n", dirname, err)
	}

	files, err := d.Readdir(-1)
	if err != nil {
		println("can't read files in folder %s: %w\n", dirname, err)
	}

	// Parse each file as a LogSet
	alllogs := LogSet{}
	for _, file := range files {
		filename := dirname + "/" + file.Name()
		data, err := readLogFile(filename)
		if err != nil {
			println("can't parse file %s: %w\n", filename, err)
		}
		for _, d := range data {
			alllogs = append(alllogs, d)
		}
	}

	// Marshal as json
	jsonLogs, err := json.Marshal(alllogs)
	if err != nil {
		println("error marshaling LogSet in JSON: %w", err)
	}

	// Output
	generateOutputs()

	os.Stdout.Write(jsonLogs)
}
