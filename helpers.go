package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func saveConfig(path string) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("can't open config file " + path + ": " + err.Error())
		return
	}
	yaml.NewEncoder(f).Encode(config)
}
