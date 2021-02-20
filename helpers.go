package main

import (
	"encoding/json"
	"log"
	"os"
)

func saveConfig(path string) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("can't open config file " + path + ": " + err.Error())
		return
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	enc.Encode(config)
}
