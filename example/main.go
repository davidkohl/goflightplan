package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/davidkohl/goflightplan"
)

func main() {
	files, err := os.ReadDir("adexp/test")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".txt" {
			filePath := filepath.Join("adexp/test", file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatal(err)
			}
			input := string(content)
			var a goflightplan.Flightplan

			err = a.Write(input)
			if err != nil {
				log.Println(err)
				continue
			}

			j, err := json.Marshal(a)
			if err != nil {
				fmt.Println(err)
			}
			_ = j
			log.Println(string(j))
		}
	}

}
