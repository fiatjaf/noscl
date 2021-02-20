package main

import (
	"log"

	"github.com/fiatjaf/go-nostr/relaypool"
)

var pool *relaypool.RelayPool

func initNostr() {
	pool = relaypool.New()

	for relay, policy := range config.Relays {
		err := pool.Add(relay, &relaypool.Policy{
			SimplePolicy: relaypool.SimplePolicy{Read: policy.Read, Write: policy.Write},
		})
		if err != nil {
			log.Printf("error adding relay '%s': %s", relay, err.Error())
		}
	}

	if len(pool.Relays) == 0 {
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
