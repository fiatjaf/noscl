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

	tags := []nostr.Tag{}
	receiverID := opts["<id>"].(string)
	tags = append(tags, nostr.Tag([]interface{}{"p", receiverID}))

	references, err := optSlice(opts, "--reference")
	if err != nil {
		return
	}
	for _, ref := range references {
		tags = append(tags, nostr.Tag([]interface{}{"e", ref}))
	}

	// parse and encrypt content
	message := opts["<content>"].(string)
	sharedSecret, err := nip04.ComputeSharedSecret(config.PrivateKey, receiverID)
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
		CreatedAt: uint32(time.Now().Unix()),
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
