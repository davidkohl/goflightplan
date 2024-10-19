package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/davidkohl/goflightplan/adexp"
)

func main() {

	logLevel := flag.String("loglevel", "info", "Log level (debug, info, warn, error)")
	dir := flag.String("dir", "", "input files to parse")
	flag.Parse()

	var level slog.Level
	switch *logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		fmt.Fprintf(os.Stderr, "invalid log level: %s\n", *logLevel)
		os.Exit(1)
	}

	Logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level, AddSource: false}))
	set, err := adexp.MessageSetFromJSON("./test/schema", "custom")
	if err != nil {
		fmt.Println(err)
		return
	}
	//m := base.MessageSet
	//m1 := adexp.MessageSet{Name: "custom", Set: base.MessageSet}
	//opts := adexp.ParserOpts{AFTNHeader: true}

	p := adexp.NewParser([]adexp.MessageSet{*set})
	_ = p

	files, err := os.ReadDir(*dir)
	if err != nil {
		fmt.Printf("could not read directory: %v\n", err)
	}
	for _, file := range files {
		// Check if it's a regular file (not a directory)
		if file.Type().IsRegular() {
			// Get the full path of the file
			filePath := *dir + "/" + file.Name()
			fmt.Println("NOW PARSING:", file.Name())
			// Read the file's contents
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("could not read file %s: %v", filePath, err)
			}
			//Print the contents of the file
			fpl, err := p.Parse(string(content))
			if err != nil {
				Logger.Error(err.Error())
				continue
			}

			j, err := json.MarshalIndent(fpl, "", "\t")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%v\n\n", string(j))

		}
	}
}
