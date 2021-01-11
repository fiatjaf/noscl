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

const USAGE = `noscli

Usage:
  noscli publish <content>
  noscli metadata --name=<name> --description=<description> --image=<image>
  noscli home [--page=<page>]
  noscli key <key> [--page=<page>]
  noscli key <key> follow
  noscli key <key> unfollow
  noscli event <id> [--page=<page>]
  noscli event <id> reference <content>
  noscli relay
  noscli relay add <url>
  noscli relay remove <url>
  noscli relay recommend <url>
`

var config struct {
	DataDir   string   `yaml:"-"`
	Relays    []Relay  `yaml:"relays,flow"`
	Following []Follow `yaml:"following,flow"`
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
	flag.StringVar(&config.DataDir, "datadir", "~/.config/nostr", "Base directory for configurations and data from Nostr.")
	flag.Parse()
	config.DataDir, _ = homedir.Expand(config.DataDir)

	// parse config
	path := filepath.Join(config.DataDir, "config.yaml")
	f, err := os.Open(path)
	if err != nil {
		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal("can't open config file " + path + ": " + err.Error())
			return
		}
		yaml.NewEncoder(f).Encode(config)
		f, _ = os.Open(path)
	}
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
	case opts["publish"].(bool):
	case opts["metadata"].(bool):
	case opts["key"].(bool):
	case opts["key"].(bool) && opts["follow"].(bool):
	case opts["key"].(bool) && opts["unfollow"].(bool):
	case opts["event"].(bool):
	case opts["event"].(bool) && opts["reference"].(bool):
	case opts["relay"].(bool):
	case opts["relay"].(bool) && opts["add"].(bool):
	case opts["relay"].(bool) && opts["remove"].(bool):
	case opts["relay"].(bool) && opts["recommend"].(bool):
	}
}
