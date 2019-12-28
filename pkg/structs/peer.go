package structs

type Peer struct {
	Hostname string `json:"hostname,omitempty"`
	Addr     string `json:"addr"`
	Port     uint16 `json:"port,omitempty"`
}

type Payload struct {
	Self  Peer   `json:"self"`
	Peers []Peer `json:"peers"`
}
