package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/davidkohl/goflightplan/adexp/models"
	"gitlab.com/davidkohl/goflightplan/adexp/schema"
)

func main() {
	v := os.Getenv("schema")
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
			var a schema.ADEXPModel
			switch v {
			case "icasv01":
				a = &models.IcasV01Model{}
			case "icasv02":
				a = &models.IcasV02Model{}
			}
			err = a.Write(input)
			if err != nil {
				fmt.Println(err)
				continue
			}

			j, err := json.Marshal(a)
			if err != nil {
				fmt.Println(err)
			}
			_ = j
			fmt.Println(string(j))
		}
	}

}
