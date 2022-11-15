package main

import (
	"log"

	"github.com/nbd-wtf/go-nostr"
)

var pool *nostr.RelayPool

func initNostr() {
	pool = nostr.NewRelayPool()

	for relay, policy := range config.Relays {
		cherr := pool.Add(relay, nostr.SimplePolicy{
			Read:  policy.Read,
			Write: policy.Write,
		})
		err := <-cherr
		if err != nil {
			log.Printf("error adding relay '%s': %s", relay, err.Error())
		}
	}

	hasRelays := false
	pool.Relays.Range(func(_ string, _ *nostr.Relay) bool {
		hasRelays = true
		return false
	})
	if !hasRelays {
		log.Printf("You have zero relays configured, everything will probably fail.")
	}

	go func() {
		for notice := range pool.Notices {
			log.Printf("%s has sent a notice: '%s'\n", notice.Relay, notice.Message)
		}
	}()

	if config.PrivateKey != "" {
		pool.SecretKey = &config.PrivateKey
	}
}
