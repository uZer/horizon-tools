package main

import (
	"encoding/json"
	"os"
)

func main() {
	// Read data folder
	dirname := "data"

	d, err := os.Open(dirname)
	if err != nil {
		println("can't open folder %s: %w\n", dirname, err)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		println("can't read files in folder %s: %w\n", dirname, err)
	}

	// Create an empty LogSet
	alllogs := LogSet{}

	// Import each file in the logset
	for _, file := range files {
		filename := dirname + "/" + file.Name()
		err := alllogs.ImportFile(filename)
		if err != nil {
			println("can't parse file %s: %w\n", filename, err)
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

	jsonTags, _ := json.Marshal(alllogs.Tags())
	os.Stdout.Write(jsonTags)

	jsonSub, _ := json.Marshal(
		alllogs.Filter(HasTag, []string{"research"}))
	os.Stdout.Write(jsonSub)
}
