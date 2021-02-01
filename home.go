package main

import (
	"github.com/docopt/docopt-go"
)

func home(opts docopt.Opts) {
	initNostr()

	for _, follow := range config.Following {
		pool.SubKey(follow.Key)
	}

	pool.ReqFeed(nil)
	printIncomingNotes()
}
