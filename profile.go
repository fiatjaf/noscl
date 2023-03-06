package main

import (
	"fmt"
	"log"

	"github.com/docopt/docopt-go"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func showProfile(opts docopt.Opts) {
	verbose, _ := opts.Bool("--verbose")
    jsonformat, _ := opts.Bool("--json")
	key := nip19.TranslatePublicKey(opts["<pubkey>"].(string))
	if key == "" {
		log.Println("Profile key is empty! Exiting.")
		return
	}

	initNostr()

	_, all := pool.Sub(nostr.Filters{{Authors: []string{key}, Kinds: []int{0}}})
	for event := range nostr.Unique(all) {
		printEvent(event, nil, verbose, jsonformat)
	}
}

func follow(opts docopt.Opts) {
	key := nip19.TranslatePublicKey(opts["<pubkey>"].(string))
	if key == "" {
		log.Println("Follow key is empty! Exiting.")
		return
	}

	name, err := opts.String("--name")
	if err != nil {
		name = ""
	}

        config.Following[key] = Follow{
		Key:  key,
		Name: name,
	}
	fmt.Printf("Followed %s.\n", key)
}

func unfollow(opts docopt.Opts) {
	key := nip19.TranslatePublicKey(opts["<pubkey>"].(string))
	if key == "" {
		log.Println("No unfollow key provided! Exiting.")
		return
	}

	delete(config.Following, key)
	fmt.Printf("Unfollowed %s.\n", key)
}

func following(opts docopt.Opts) {
	if len(config.Following) == 0 {
		fmt.Println("You aren't following anyone yet.")
		return
	}
	for _, profile := range config.Following {
		fmt.Println(profile.Key, profile.Name)
	}
}
