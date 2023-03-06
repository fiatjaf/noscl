package main

import (
	"log"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
)

func message(opts docopt.Opts) {
	if config.PrivateKey == "" {
		log.Printf("Can't direct message. Private key not set.\n")
		return
	}

	initNostr()

	var tags nostr.Tags
	receiverKey := opts["<pubkey>"].(string)
	tags = append(tags, nostr.Tag{"p", receiverKey})

	references, err := optSlice(opts, "--reference")
	if err != nil {
		return
	}
	for _, ref := range references {
		tags = append(tags, nostr.Tag{"e", ref})
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

    log.Printf("%+v\n", event)
    log.Printf("%+v\n", statuses)
	// printPublishStatus(event, statuses)
}
