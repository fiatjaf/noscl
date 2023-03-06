package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"gopkg.in/yaml.v2"
)

var kindNames = map[int]string{
	nostr.KindSetMetadata:            "Profile Metadata",
	nostr.KindTextNote:               "Text Note",
	nostr.KindRecommendServer:        "Relay Recommendation",
	nostr.KindContactList:            "Contact List",
	nostr.KindEncryptedDirectMessage: "Encrypted Message",
	nostr.KindDeletion:               "Deletion Notice",
	nostr.KindBoost:                  "Boost",
	nostr.KindReaction:               "Reaction",
	nostr.KindChannelCreation:        "Channel Creation",
	nostr.KindChannelMetadata:        "Channel Metadata",
	nostr.KindChannelMessage:         "Channel Message",
	nostr.KindChannelHideMessage:     "Channel Hide Message",
	nostr.KindChannelMuteUser:        "Channel Mute User",
}

func printEvent(evt nostr.Event, nick *string, verbose bool, jsonformat bool) {
	kind, ok := kindNames[evt.Kind]
	if !ok {
		kind = fmt.Sprintf("Unknown Kind (%d)", evt.Kind)
	}

	// Don't print encrypted messages that aren't for me or from me
    pubkey := getPubKey(config.PrivateKey)
	if evt.Kind == nostr.KindEncryptedDirectMessage {
		if (!evt.Tags.ContainsAny("p", nostr.Tag{getPubKey(config.PrivateKey)})) && (evt.PubKey != pubkey) {
			return
		}
	}

    // json
    if jsonformat {
        jevt, _ := json.MarshalIndent(evt, "", "\t")
        fmt.Print(string(jevt))
        return
    }

	var ID string = shorten(evt.ID)
	var fromField string = shorten(evt.PubKey)

	if nick != nil {
		fromField = fmt.Sprintf("%s (%s)", *nick, shorten(evt.PubKey))
	}

	if verbose {
		ID = evt.ID

		if nick == nil {
			fromField = evt.PubKey
		} else {
			fromField = fmt.Sprintf("%s (%s)", *nick, evt.PubKey)
		}
	}

	fmt.Printf("%s [%s] from %s %s\n",
		kind,
		ID,
		fromField,
		humanize.Time(evt.CreatedAt),
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
		spl := strings.Split(string(y), "\n")
		for i, v := range spl {
			spl[i] = "  " + v
		}
		str := strings.Join(spl, "\n")
		fmt.Print(str)
	case nostr.KindTextNote:
		fmt.Print("  " + strings.ReplaceAll(evt.Content, "\n", "\n  "))
	case nostr.KindBoost:
		var event nostr.Event
		err := json.Unmarshal([]byte(evt.Content), &event)
		if err != nil {
			fmt.Println("ERR", err)
			return
		}

		kind, ok := kindNames[event.Kind]
		if !ok {
			kind = "Unknown Kind"
		}
		var ID string = shorten(event.ID)
		var fromField string = shorten(event.PubKey)
		fmt.Printf("  %s [%s] from %s %s\n",
			kind,
			ID,
			fromField,
			humanize.Time(event.CreatedAt),
		)
		fmt.Print("    ", event.Content)
	case nostr.KindRecommendServer:
	case nostr.KindContactList:
	case nostr.KindEncryptedDirectMessage:
		sharedSecret, err := nip04.ComputeSharedSecret(config.PrivateKey, evt.PubKey)
		if err != nil {
			log.Printf("Error computing shared key: %s. \n", err.Error())
			return
		}
		txt, err := nip04.Decrypt(evt.Content, sharedSecret)
		if err != nil {
			log.Printf("Error decrypting message: %s. \n", err.Error())
			return
		}
		fmt.Print(txt)

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
