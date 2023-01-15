package main

var config Config

type Config struct {
	DataDir    string            `json:"-"`
	Relays     map[string]Policy `json:"relays,flow"`
	Following  map[string]Follow `json:"following,flow"`
	PrivateKey string            `json:"privatekey,omitempty"`
}

type Follow struct {
	Key    string   `json:"key"`
	Name   string   `json:"name,flow,omitempty"`
	Relays []string `json:"relays,flow,omitempty"`
}

type Policy struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
}

func (p Policy) String() string {
	var ret string
	if p.Read {
		ret += "r"
	}
	if p.Write {
		ret += "w"
	}
	return ret
}

func (c *Config) Init() {
	if c.Relays == nil {
		c.Relays = make(map[string]Policy)
	}
	if c.Following == nil {
		c.Following = make(map[string]Follow)
	}
}
