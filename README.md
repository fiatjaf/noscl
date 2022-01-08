noscl
=====

Command line client for [Nostr](https://github.com/fiatjaf/nostr).

## Installation

Compile with `go get github.com/fiatjaf/noscl` or [download a binary](releases/).

## Usage

```
noscl

Usage:
  noscl home
  noscl setprivate <key>
  noscl public
  noscl publish [--reference=<id>] <content>
  noscl metadata --name=<name> [--description=<description>] [--image=<image>]
  noscl profile <key>
  noscl follow <key>
  noscl unfollow <key>
  noscl event <id>
  noscl share-contants
  noscl key-gen
  noscl relay
  noscl relay add <url>
  noscl relay remove <url>
  noscl relay recommend <url>
```

The basic flow is something like

1. Add some relays with `noscl relay add <relay url>` (see https://nostr-registry.netlify.app/ for some publicly available relays)
2. Follow some people with `noscl follow <pubkey>`
3. Browse some profiles from people (you don't have to be following) with `noscl profile <pubkey>`
4. See the feed of people you follow with `noscl home`
5. Set your own private key with `noscl setprivate <hex private key>`
6. Get your public key with `noscl public` so you can share it with others
7. Publish some notes with `noscl publish <my note content>`

## Generate a key

```
$ noscl key-gen
seed: crowd coconut donate calm position chuckle update friend ball gospel sudden answer bitter dinosaur wise express jaguar file praise pact defy system struggle offer
private key: 5a860fa953c9162611f2e2813ee0526370664534412f31611a4a89149c6bbc9d

$ noscl setprivate 5a860fa953c9162611f2e2813ee0526370664534412f31611a4a89149c6bbc9d
```
