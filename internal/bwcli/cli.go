package bwcli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

const defaultAddress = "127.0.0.1:58770"

func Run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		printUsage(stdout)
		return nil
	}

	address := defaultAddress
	if args[0] == "--addr" {
		if len(args) < 3 {
			return errors.New("--addr requires an address and a command")
		}
		address = args[1]
		args = args[2:]
	}

	cmd, err := commandFromArgs(args)
	if err != nil {
		return err
	}
	if cmd.Action == "help" {
		printUsage(stdout)
		return nil
	}

	reply, err := send(ctx, address, cmd)
	if err != nil {
		return err
	}
	if reply != "" {
		fmt.Fprintln(stdout, reply)
	}
	return nil
}

func commandFromArgs(args []string) (Command, error) {
	switch args[0] {
	case "help", "-h", "--help":
		return Command{Action: "help"}, nil
	case "status":
		return Command{Action: "status"}, nil
	case "play", "stop", "record", "overdub", "panic":
		return Command{Action: args[0]}, nil
	case "tempo":
		if len(args) != 2 {
			return Command{}, errors.New("usage: bw tempo <bpm>")
		}
		bpm, err := strconv.ParseFloat(args[1], 64)
		if err != nil || bpm <= 0 {
			return Command{}, fmt.Errorf("invalid tempo %q", args[1])
		}
		return Command{Action: "tempo", Value: bpm}, nil
	case "track":
		if len(args) >= 2 && args[1] == "create" {
			if len(args) != 3 && len(args) != 4 {
				return Command{}, errors.New("usage: bw track create <instrument|audio|effect> [position]")
			}
			cmd := Command{Action: "track.create." + args[2]}
			if len(args) == 4 {
				number, err := strconv.Atoi(args[3])
				if err != nil || number < 1 {
					return Command{}, fmt.Errorf("invalid track position %q", args[3])
				}
				cmd.Track = number
			}
			return cmd, nil
		}
		if len(args) != 3 {
			return Command{}, errors.New("usage: bw track <number> <arm|mute|solo|select>")
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return Command{}, fmt.Errorf("invalid track number %q", args[1])
		}
		return Command{Action: "track." + args[2], Track: number}, nil
	case "scene":
		if len(args) == 2 && args[1] == "create" {
			return Command{Action: "scene.create"}, nil
		}
		if len(args) != 3 {
			return Command{}, errors.New("usage: bw scene create | bw scene <number> <select|show|launch>")
		}
		number, err := strconv.Atoi(args[1])
		if err != nil || number < 1 {
			return Command{}, fmt.Errorf("invalid scene number %q", args[1])
		}
		return Command{Action: "scene." + args[2], Slot: number}, nil
	case "clip":
		if len(args) < 4 {
			return Command{}, errors.New("usage: bw clip <track> <slot> <launch|record|stop|select|show|duplicate|delete|clear|create|note|pattern>")
		}
		track, err := strconv.Atoi(args[1])
		if err != nil || track < 1 {
			return Command{}, fmt.Errorf("invalid track number %q", args[1])
		}
		slot, err := strconv.Atoi(args[2])
		if err != nil || slot < 1 {
			return Command{}, fmt.Errorf("invalid clip slot %q", args[2])
		}
		if args[3] == "create" {
			if len(args) != 5 {
				return Command{}, errors.New("usage: bw clip <track> <slot> create <beats>")
			}
			beats, err := strconv.ParseFloat(args[4], 64)
			if err != nil || beats <= 0 {
				return Command{}, fmt.Errorf("invalid clip length %q", args[4])
			}
			return Command{Action: "clip.create", Track: track, Slot: slot, Value: beats}, nil
		}
		if args[3] == "note" {
			if len(args) != 7 {
				return Command{}, errors.New("usage: bw clip <track> <slot> note <pitch> <step> <length>")
			}
			pitch, err := strconv.Atoi(args[4])
			if err != nil || pitch < 0 || pitch > 127 {
				return Command{}, fmt.Errorf("invalid MIDI pitch %q", args[4])
			}
			step, err := strconv.Atoi(args[5])
			if err != nil || step < 0 {
				return Command{}, fmt.Errorf("invalid step %q", args[5])
			}
			length, err := strconv.ParseFloat(args[6], 64)
			if err != nil || length <= 0 {
				return Command{}, fmt.Errorf("invalid note length %q", args[6])
			}
			return Command{Action: "clip.note", Track: track, Slot: slot, Pitch: pitch, Step: step, Length: length}, nil
		}
		if args[3] == "pattern" {
			if len(args) != 5 {
				return Command{}, errors.New("usage: bw clip <track> <slot> pattern <name>")
			}
			return Command{Action: "clip.pattern", Track: track, Slot: slot, Name: args[4]}, nil
		}
		if len(args) != 4 {
			return Command{}, errors.New("usage: bw clip <track> <slot> <launch|record|stop|select|show|duplicate|delete|clear>")
		}
		return Command{Action: "clip." + args[3], Track: track, Slot: slot}, nil
	default:
		return Command{}, fmt.Errorf("unknown command %q", args[0])
	}
}

func send(ctx context.Context, address string, cmd Command) (string, error) {
	payload, err := cmd.Marshal()
	if err != nil {
		return "", err
	}

	var d net.Dialer
	conn, err := d.DialContext(ctx, "udp", address)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}

	if _, err := conn.Write(payload); err != nil {
		return "", err
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("no response from Bitwig bridge at %s: %w", address, err)
	}
	return string(buf[:n]), nil
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, `bw controls Bitwig through the local Bitwig CLI Bridge extension.

Usage:
  bw [--addr 127.0.0.1:58770] status
  bw play|stop|record|overdub|panic
  bw tempo <bpm>
  bw track create <instrument|audio|effect> [position]
  bw track <number> <arm|mute|solo|select>
  bw scene create
  bw scene <number> <select|show|launch>
  bw clip <track> <slot> <launch|record|stop|select|show|duplicate|delete|clear>
  bw clip <track> <slot> create <beats>
  bw clip <track> <slot> note <pitch> <step> <length>
  bw clip <track> <slot> pattern <test-c|bass-c>`)
}
