package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/fiatjaf/go-nostr"
	"gopkg.in/yaml.v2"
	"math/big"
	"strings"
	"time"
)

var kindNames = map[int]string{
	nostr.KindSetMetadata:            "Profile Metadata",
	nostr.KindTextNote:               "Text Note",
	nostr.KindRecommendServer:        "Relay Recommendation",
	nostr.KindContactList:            "Contact List",
	nostr.KindEncryptedDirectMessage: "Encrypted Message",
	nostr.KindDeletion:               "Deletion Notice",
	nostr.KindPow:                    "Proof of Work",
}

func printEvent(evt nostr.Event, nick *string) {
	kind, ok := kindNames[evt.Kind]
	if !ok {
		kind = "Unknown Kind"
	}

	// Don't print encrypted messages that aren't for me
	if evt.Kind == nostr.KindEncryptedDirectMessage {
		if !evt.Tags.ContainsAny("p", nostr.StringList{getPubKey(config.PrivateKey)}) {
			return
		}
	}

	var fromField string

	if nick == nil {
		fromField = shorten(evt.PubKey)
	} else {
		fromField = fmt.Sprintf("%s (%s)", *nick, shorten(evt.PubKey))
	}

	fmt.Printf("%s [%s] from %s %s\n",
		kind,
		shorten(evt.ID),
		fromField,
		humanize.Time(time.Unix(int64(evt.CreatedAt), 0)),
	)

	switch evt.Kind {
	case nostr.KindSetMetadata:
		var metadata Metadata
		err := json.Unmarshal([]byte(evt.Content), &metadata)
		if err != nil {
			fmt.Printf("Invalid JSON: '%s',\n  %s",
				err.Error(), evt.Content)
			return
		}
		y, _ := yaml.Marshal(metadata)
		fmt.Print(string(y))
	case nostr.KindTextNote:
		if pows, ok := checkPow(evt); ok {
			for alg, pow := range pows {
				fmt.Printf(" ↑ %s %s\n", alg, shortPowHash(pow.hash))
			}
		}
		fmt.Printf("  ")
		fmt.Print(strings.ReplaceAll(evt.Content, "\n", "\n  "))
	case nostr.KindRecommendServer:
	case nostr.KindContactList:
	case nostr.KindEncryptedDirectMessage:
	case nostr.KindPow:
		if pows, ok := checkPowVote(evt); ok {
			for alg, pow := range pows {
				fmt.Printf(" ↑ %s %s", alg, shortPowHash(pow.hash))
			}
		}
	default:
		fmt.Print(evt.Content)
	}

	fmt.Printf("\n")
}

func shorten(id string) string {
	if len(id) < 12 {
		return id
	}
	return id[0:4] + "..." + id[len(id)-4:]
}

// A hash is an integer with a maximum value of MAX = 256^(len(hash))
// The probability that a single hash is less than some number x is:
// 	p = x / MAX
// The number of hashes, k, needed to generate a hash less than x,
// follows a geometric distribution (wikipedia.org/wiki/Geometric_distribution):
// 	Pr(X=k) = (1-p)^{k-1} p
// Thus, for a given x, the expected number of hashes is:
// 	E(X) = 1/p = MAX/x
func shortPowHash(hash []byte) string {
	var x, MAX, r big.Int
	x.SetBytes(hash)
	MAX.Exp(big.NewInt(256), big.NewInt(int64(len(hash))), nil)
	r.Div(&MAX, &x)
	return fmt.Sprintf("+%s (%s...)", r.String(), hex.EncodeToString(hash)[0:12])
}

func printPublishStatus(event *nostr.Event, statuses chan nostr.PublishStatus) {
	for status := range statuses {
		switch status.Status {
		case nostr.PublishStatusSent:
			fmt.Printf("Sent event %s to '%s'.\n", event.ID, status.Relay)
		case nostr.PublishStatusFailed:
			fmt.Printf("Failed to send event %s to '%s'.\n", event.ID, status.Relay)
		case nostr.PublishStatusSucceeded:
			fmt.Printf("Seen %s on '%s'.\n", event.ID, status.Relay)
		}
	}
}
