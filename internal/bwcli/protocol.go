package bwcli

import "encoding/json"

type Command struct {
	Action string  `json:"action"`
	Track  int     `json:"track,omitempty"`
	Slot   int     `json:"slot,omitempty"`
	Pitch  int     `json:"pitch,omitempty"`
	Step   int     `json:"step,omitempty"`
	Length float64 `json:"length,omitempty"`
	Value  float64 `json:"value,omitempty"`
	Name   string  `json:"name,omitempty"`
}

func (c Command) Marshal() ([]byte, error) {
	return json.Marshal(c)
}
