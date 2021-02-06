package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec"
	"github.com/docopt/docopt-go"
)

func setPrivateKey(opts docopt.Opts) {
	keyhex := opts["<key>"].(string)

	if len(keyhex) < 64 {
		log.Print("key too short, must be 32 bytes hex-encoded, i.e. 64 characters.\n")
		return
	}

	if _, err := hex.DecodeString(keyhex); err != nil {
		log.Printf("Error decoding key from hex: %s\n", err.Error())
		return
	}

	config.PrivateKey = keyhex
}

func showPublicKey(opts docopt.Opts) {
	if config.PrivateKey == "" {
		log.Printf("No private key set.\n")
		return
	}

	pubkey := getPubKey(config.PrivateKey)
	if pubkey != "" {
		fmt.Printf("%x\n", pubkey)
	}
}

func getPubKey(privateKey string) string {
	if keyb, err := hex.DecodeString(config.PrivateKey); err != nil {
		log.Printf("Error decoding key from hex: %s\n", err.Error())
		return ""
	} else {
		_, pubkey := btcec.PrivKeyFromBytes(btcec.S256(), keyb)
		return hex.EncodeToString(pubkey.SerializeCompressed()[1:])
	}
}
