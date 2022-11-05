package main

import (
	"log"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/fiatjaf/go-nostr"
	"github.com/fiatjaf/go-nostr/nip04"
)

func message(opts docopt.Opts) {
	if config.PrivateKey == "" {
		log.Printf("Can't direct message. Private key not set.\n")
		return
	}

	initNostr()

	var tags nostr.Tags
	receiverKey := opts["<pubkey>"].(string)
	tags = append(tags, nostr.StringList{"p", receiverKey})

	references, err := optSlice(opts, "--reference")
	if err != nil {
		return
	}
	for _, ref := range references {
		tags = append(tags, nostr.StringList{"e", ref})
	}

	// parse and encrypt content
	message := opts["<content>"].(string)
	if message == "-" {
		message, err = readContentStdin(4096)
		if err != nil {
			log.Printf("Failed reading content from stdin: %s", err)
		}
	}
	sharedSecret, err := nip04.ComputeSharedSecret(config.PrivateKey, receiverKey)
	if err != nil {
		log.Printf("Error computing shared key: %s. \n", err.Error())
		return
	}

	encryptedMessage, err := nip04.Encrypt(message, sharedSecret)
	if err != nil {
		log.Printf("Error encrypting message: %s. \n", err.Error())
		return
	}

	event, statuses, err := pool.PublishEvent(&nostr.Event{
		CreatedAt: time.Now(),
		Kind:      nostr.KindEncryptedDirectMessage,
		Tags:      tags,
		Content:   encryptedMessage,
	})
	if err != nil {
		log.Printf("Error messaging: %s.\n", err.Error())
		return
	}

	printPublishStatus(event, statuses)
}
