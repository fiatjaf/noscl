package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/docopt/docopt-go"
	"github.com/nbd-wtf/go-nostr"
)

func signEventJSON(opts docopt.Opts) {
	if config.PrivateKey == "" {
		log.Printf("Can't sign. Private key not set.\n")
		return
	}

	j, err := opts.String("<event-json>")
	if err != nil {
		log.Printf("Provide the JSON event as a quoted argument: %s.\n", err.Error())
		return
	}

	var event nostr.Event
	if err := json.Unmarshal([]byte(j), &event); err != nil {
		log.Printf("Invalid event JSON: %s.\n", err.Error())
		return
	}

	if err := event.Sign(config.PrivateKey); err != nil {
		log.Printf("Failed to sign: %s.\n", err.Error())
		return
	}

	fmt.Printf("Id: %s\nSignature: %s\n", event.ID, event.Sig)
}

func verifyEventJSON(opts docopt.Opts) {
	j, err := opts.String("<event-json>")
	if err != nil {
		log.Printf("Provide the JSON event as a quoted argument: %s.\n", err.Error())
		return
	}

	var event nostr.Event
	if err := json.Unmarshal([]byte(j), &event); err != nil {
		log.Printf("Invalid event JSON: %s.\n", err.Error())
		return
	}

	if ok, err := event.CheckSignature(); err != nil {
		fmt.Printf("Serialized: %s\n", event.Serialize())
		fmt.Printf("Hash: %s\n", event.GetID())
		fmt.Printf("Failed to verify: %s.\n", err.Error())
		return
	} else if !ok {
		fmt.Printf("Signature is invalid.\n")
		return
	}

	fmt.Printf("Signature is valid.\n")
}
