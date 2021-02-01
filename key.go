package main

import (
	"fmt"

	"github.com/docopt/docopt-go"
)

func showKey(opts docopt.Opts) {
	key := opts["<key>"].(string)
	pool.ReqKey(key, nil)
	printIncomingNotes()
}

func follow(opts docopt.Opts) {
	key := opts["<key>"].(string)
	config.Following = append(config.Following, Follow{
		Key: key,
	})
	fmt.Printf("Followed %s.\n", key)
}

func unfollow(opts docopt.Opts) {
	key := opts["<key>"].(string)
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
