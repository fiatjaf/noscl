package main

import (
	"encoding/json"
	"fmt"
	"io"
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

// readContentStdin reads from stdin until EOF or up to max + 1 bytes.
// it return an error if the read length is larger than max.
func readContentStdin(max int) (string, error) {
	b, err := io.ReadAll(io.LimitReader(os.Stdin, int64(max)+1))
	if err != nil {
		return "", err
	}
	if len(b) == max+1 {
		return "", fmt.Errorf("too big; want max %d bytes", max)
	}
	return string(b), nil
}
