package main

import (
	"fmt"

	"github.com/docopt/docopt-go"
)

func addRelay(opts docopt.Opts) {
	addr := opts["<url>"].(string)
	config.Relays = append(config.Relays, Relay{addr, "rw"})
	fmt.Printf("Added relay %s.\n", addr)
}

func removeRelay(opts docopt.Opts) {
	addr := opts["<url>"].(string)
	var newRelaysList []Relay
	for _, relay := range config.Relays {
		if relay.URL == addr {
			continue
		}
		newRelaysList = append(newRelaysList, relay)
	}
	config.Relays = newRelaysList
	fmt.Printf("Removed relay %s.\n", addr)
}

func recommendRelay(opts docopt.Opts) {
	addr := opts["<url>"].(string)

	// TODO

	fmt.Printf("Published a relay recommendation for %s.", addr)
}

func listRelays(opts docopt.Opts) {
	for _, relay := range config.Relays {
		fmt.Printf("%s: %s\n", relay.URL, relay.Policy)
	}
}
