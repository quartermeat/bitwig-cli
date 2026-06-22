package bwcli

import "testing"

func TestCommandFromArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want Command
	}{
		{name: "play", args: []string{"play"}, want: Command{Action: "play"}},
		{name: "tempo", args: []string{"tempo", "128"}, want: Command{Action: "tempo", Value: 128}},
		{name: "track create", args: []string{"track", "create", "instrument", "8"}, want: Command{Action: "track.create.instrument", Track: 8}},
		{name: "track arm", args: []string{"track", "2", "arm"}, want: Command{Action: "track.arm", Track: 2}},
		{name: "scene create", args: []string{"scene", "create"}, want: Command{Action: "scene.create"}},
		{name: "scene select", args: []string{"scene", "3", "select"}, want: Command{Action: "scene.select", Slot: 3}},
		{name: "clip launch", args: []string{"clip", "2", "4", "launch"}, want: Command{Action: "clip.launch", Track: 2, Slot: 4}},
		{name: "clip create", args: []string{"clip", "2", "4", "create", "8"}, want: Command{Action: "clip.create", Track: 2, Slot: 4, Value: 8}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := commandFromArgs(tt.args)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCommandFromArgsRejectsBadTempo(t *testing.T) {
	if _, err := commandFromArgs([]string{"tempo", "fast"}); err == nil {
		t.Fatal("expected error")
	}
}
