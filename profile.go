package main

import (
	"fmt"
	"log"

	"github.com/docopt/docopt-go"
	"github.com/fiatjaf/go-nostr"
)

func showProfile(opts docopt.Opts) {
	initNostr()

	key := opts["<key>"].(string)
	if key == "" {
		log.Println("Profile key is empty! Exiting.")
		return
	}

	sub := pool.Sub(nostr.EventFilters{{Authors: []string{key}}})
	for event := range sub.UniqueEvents {
		printEvent(event)
	}
}

func follow(opts docopt.Opts) {
	key := opts["<key>"].(string)
	if key == "" {
		log.Println("Follow key is empty! Exiting.")
		return
	}

	name, err := opts.String("--name")

	if err != nil {
		name = ""
	}

	config.Following = append(config.Following, Follow{
		Key:  key,
		Name: name,
	})
	fmt.Printf("Followed %s.\n", key)
}

func unfollow(opts docopt.Opts) {
	key := opts["<key>"].(string)
	if key == "" {
		log.Println("No unfollow key provided! Exiting.")
		return
	}

	var newFollowingList []Follow
	for _, follow := range config.Following {
		if follow.Key == key {
			continue
		}
		newFollowingList = append(newFollowingList, follow)
	}
	config.Following = newFollowingList
	fmt.Printf("Unfollowed %s.\n", key)
}
