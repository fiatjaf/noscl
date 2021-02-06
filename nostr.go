package main

import (
	"fmt"
	"os"

	"github.com/fiatjaf/go-nostr/relaypool"
)

var pool *relaypool.RelayPool

func initNostr() {
	pool = relaypool.New()

	for _, relay := range config.Relays {
		pool.Add(relay.URL, nil)
	}

	for notice := range pool.Notices {
		fmt.Fprintf(os.Stderr, "%s sent a notice: '%s'\n", notice.Relay, notice.Message)
	}

	if config.PrivateKey != "" {
		pool.SecretKey = &config.PrivateKey
	}
}
