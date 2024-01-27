noscl
=====
[![Go Report Card](https://goreportcard.com/badge/github.com/fiatjaf/noscl)](https://goreportcard.com/report/github.com/fiatjaf/noscl)  [![License: ODbL](https://img.shields.io/badge/License-PDDL-brightgreen.svg)](https://opendatacommons.org/licenses/pddl/)  [![Latest Release](https://img.shields.io/github/v/release/fiatjaf/noscl?logo=github)](https://github.com/fiatjaf/noscl/releases)

Command line client for [Nostr](https://github.com/fiatjaf/nostr).

## Notice

Although it works, this project is somewhat abandoned. For a more complete CLI experience you can try https://github.com/mattn/algia and for a more streamlined CLI plumbing tool try https://github.com/fiatjaf/nak.

## Installation

Compile with `go install github.com/fiatjaf/noscl@latest` or [download a binary](https://github.com/fiatjaf/noscl/releases).

## Usage

```
noscl

Usage:
  noscl home
  noscl setprivate <key>
  noscl public
  noscl publish [--reference=<id>...] [--profile=<id>...] <content>
  noscl message [--reference=<id>...] <id> <content>
  noscl metadata --name=<name> [--about=<about>] [--picture=<picture>]
  noscl profile <key>
  noscl follow <key> [--name=<name>]
  noscl unfollow <key>
  noscl following
  noscl event <id>
  noscl share-contacts
  noscl key-gen
  noscl relay
  noscl relay add <url>
  noscl relay remove <url>
  noscl relay recommend <url>
```

The basic flow is something like

1. Add some relays with `noscl relay add <relay url>` (see https://nostr.watch/relays/find for some publicly available relays)
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

## Sign an event manually

```
noscl sign '{...event as JSON}'
```

https://user-images.githubusercontent.com/1653275/149637925-32943e2e-a2ff-41a0-9e3d-5ea1a60c84ae.mp4
