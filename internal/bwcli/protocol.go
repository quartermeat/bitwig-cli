package bwcli

import "encoding/json"

type Command struct {
	Action string  `json:"action"`
	Track  int     `json:"track,omitempty"`
	Slot   int     `json:"slot,omitempty"`
	Value  float64 `json:"value,omitempty"`
}

func (c Command) Marshal() ([]byte, error) {
	return json.Marshal(c)
}
