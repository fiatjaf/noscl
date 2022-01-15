package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/fiatjaf/go-nostr"
	"golang.org/x/crypto/scrypt"
)

type checkedPow struct {
	nonce string
	hash  []byte
}

func checkPowVote(evt nostr.Event) (map[string]checkedPow, bool) {
	var e string
	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		currentTagName, ok := tag[0].(string)
		if !ok || currentTagName != "e" {
			continue
		}
		currentTagValue, ok := tag[1].(string)
		if !ok {
			continue
		}
		e = currentTagValue
		break
	}
	if e == "" {
		log.Printf("No referenced event,\n")
		return nil, false
	}
	evt.ID = e

	pow := make([][]string, 0, 1)
	if err := json.Unmarshal([]byte(evt.Content), &pow); err != nil {
		log.Printf("Could not decode POW: %s", evt.Content)
		return nil, false
	}
	evt.Pow = pow

	return checkPow(evt)
}

func checkPow(evt nostr.Event) (map[string]checkedPow, bool) {
	id, err := hex.DecodeString(evt.ID)
	if err != nil {
		log.Printf("Could not decode id: %s.\n", err.Error())
		return nil, false
	}
	result := make(map[string]checkedPow, 1)
	for _, pow := range evt.Pow {
		switch pow[0] {
		case "scrypt":
			nonce, err := hex.DecodeString(pow[1])
			if err != nil {
				log.Printf("Could not decode nonce: %s.\n", err.Error())
				return nil, false
			}
			hash, err := scrypt.Key(id, nonce, 32768, 8, 1, 32)
			if err != nil {
				log.Printf("Error hashsing scrypt: %s.\n", err.Error())
				return nil, false
			}
			if opow, ok := result["scrypt"]; ok {
				if bytes.Compare(hash, opow.hash) < 0 {
					// hash less than opow.hash
					result["scrypt"] = checkedPow{pow[1], hash}
				}
			} else {
				result["scrypt"] = checkedPow{pow[1], hash}
			}
		}
	}
	return result, true
}

func powScrypt(message []byte, n int) ([]byte, error) {

	bestnonce := make([]byte, 8)
	nonce := make([]byte, 8)
	best := []byte{
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255}
	for i := 0; i < n; i++ {
		rand.Read(nonce)
		other, err := scrypt.Key(message, nonce, 32768, 8, 1, 32)
		if err != nil {
			return bestnonce, err
		}
		if bytes.Compare(other, best) < 0 { // other less than best
			best = other
			copy(bestnonce, nonce)
		}
	}

	return bestnonce, nil
}

func pow(opts docopt.Opts) {
	initNostr()
	rand.Seed(time.Now().UnixNano())

	e := opts["<id>"].(string)
	message, err := hex.DecodeString(e)
	if err != nil {
		log.Printf("Could not decode event id: %s.\n", err.Error())
		return
	}

	tags := make(nostr.Tags, 0, 1)
	tags = append(tags, nostr.Tag([]interface{}{"e", e}))

	n, err := opts.Int("<n>")
	if err != nil {
		log.Printf("Not a number of hashes to perform: %s.\n", err.Error())
		return
	}

	nonce, err := powScrypt(message, n)
	if err != nil {
		log.Printf("POW error: %s.\n", err.Error())
		return
	}

	pow, _ := json.Marshal([][]string{{"scrypt", hex.EncodeToString(nonce)}})

	event, statuses, err := pool.PublishEvent(&nostr.Event{
		CreatedAt: uint32(time.Now().Unix()),
		Kind:      nostr.KindPow,
		Tags:      tags,
		Content:   string(pow),
	})
	if err != nil {
		log.Printf("Error publishing: %s.\n", err.Error())
		return
	}

	printPublishStatus(event, statuses)
}
