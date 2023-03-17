package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/bech32"

	"github.com/docopt/docopt-go"
	"github.com/nbd-wtf/go-nostr/nip06"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func decodeKey(keyraw string) ([]byte, error) {
	if len(keyraw) == 64 {
		// hex-encoded
		keyval, err := hex.DecodeString(keyraw)
		if err != nil {
			return nil, fmt.Errorf("decoding key from hex: %w", err)
		}
		return keyval, nil
	}

	// bech32-encoded
	_, keyval, err := bech32.Decode(keyraw)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("decoding key from bech32: %w", err)
	}
	return keyval, nil
}

func setPrivateKey(opts docopt.Opts) {
	keyraw := opts["<key>"].(string)
	keyval, err := decodeKey(keyraw)
	if err != nil {
		log.Printf("Failed to parse private key: %s\n", err.Error())
		return
	}

	config.PrivateKey = string(keyval)
}

func showPublicKey(opts docopt.Opts) {
	if config.PrivateKey == "" {
		log.Printf("No private key set.\n")
		return
	}

	pubkey := getPubKey(config.PrivateKey)
	if pubkey != "" {
		fmt.Printf("%s\n", pubkey)

		nip19pubkey, _ := nip19.EncodePublicKey(pubkey, "")
		fmt.Printf("%s\n", nip19pubkey)
	}
}

func getPubKey(privateKey string) string {
	if keyb, err := hex.DecodeString(privateKey); err != nil {
		log.Printf("Error decoding key from hex: %s\n", err.Error())
		return ""
	} else {
		_, pubkey := btcec.PrivKeyFromBytes(keyb)
		return hex.EncodeToString(schnorr.SerializePubKey(pubkey))
	}
}

func keyGen(opts docopt.Opts) {
	seedWords, err := nip06.GenerateSeedWords()
	if err != nil {
		log.Println(err)
		return
	}

	seed := nip06.SeedFromWords(seedWords)

	sk, err := nip06.PrivateKeyFromSeed(seed)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("seed:", seedWords)
	fmt.Println("private key:", sk)
}
