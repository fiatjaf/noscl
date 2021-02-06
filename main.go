package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/docopt/docopt-go"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

const USAGE = `noscl

Usage:
  noscl home [--page=<page>]
  noscl setprivate <key>
  noscl public
  noscl publish [--reference=<id>] <content>
  noscl metadata --name=<name> --description=<description> --image=<image>
  noscl profile <key> [--page=<page>]
  noscl follow <key>
  noscl unfollow <key>
  noscl event <id> [--page=<page>]
  noscl relay
  noscl relay add <url>
  noscl relay remove <url>
  noscl relay recommend <url>
`

var config struct {
	DataDir    string   `yaml:"-"`
	Relays     []Relay  `yaml:"relays,flow"`
	Following  []Follow `yaml:"following,flow"`
	PrivateKey string   `yaml:"privatekey,omitempty"`
}

type Relay struct {
	URL    string `yaml:"url"`
	Policy string `yaml:"policy"` // "r" for read, "w" for write, "n" for no-related
}

type Follow struct {
	Key    string   `yaml:"key"`
	Name   string   `yaml:"name,omitempty"`
	Relays []string `yaml:"relays,omitempty"`
}

func main() {
	// find datadir
	flag.StringVar(&config.DataDir, "datadir", "~/.config/nostr",
		"Base directory for configurations and data from Nostr.")
	flag.Parse()
	config.DataDir, _ = homedir.Expand(config.DataDir)
	os.Mkdir(config.DataDir, 0700)

	// parse config
	path := filepath.Join(config.DataDir, "config.yaml")
	f, err := os.Open(path)
	if err != nil {
		saveConfig(path)
		f, _ = os.Open(path)
	}
	f, _ = os.Open(path)
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		log.Fatal("can't parse config file " + path + ": " + err.Error())
		return
	}

	// parse args
	opts, err := docopt.ParseDoc(USAGE)
	if err != nil {
		return
	}

	switch {
	case opts["home"].(bool):
		home(opts)
	case opts["setprivate"].(bool):
		// TODO make this read STDIN and encrypt the key locally
		// TODO also accept BIP39
		setPrivateKey(opts)
		saveConfig(path)
	case opts["public"].(bool):
		showPublicKey(opts)
	case opts["publish"].(bool):
		publish(opts)
	case opts["metadata"].(bool):
	case opts["profile"].(bool):
		showProfile(opts)
	case opts["follow"].(bool):
		follow(opts)
		saveConfig(path)
	case opts["unfollow"].(bool):
		unfollow(opts)
		saveConfig(path)
	case opts["event"].(bool):
		view(opts)
	case opts["relay"].(bool):
		switch {
		case opts["add"].(bool):
			addRelay(opts)
			saveConfig(path)
		case opts["remove"].(bool):
			removeRelay(opts)
			saveConfig(path)
		case opts["recommend"].(bool):
			recommendRelay(opts)
		default:
			listRelays(opts)
		}
	}
}
