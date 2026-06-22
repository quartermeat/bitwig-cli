# Bitwig CLI

`bitwig-cli` is a local Go command-line tool for controlling a live Bitwig
performance through a companion Bitwig Controller Extension.

The first bridge design is intentionally simple:

```text
bw.exe -> UDP JSON on 127.0.0.1:58770 -> Bitwig Controller Extension -> Bitwig
```

## Commands

```powershell
go run ./cmd/bw status
go run ./cmd/bw play
go run ./cmd/bw stop
go run ./cmd/bw save
go run ./cmd/bw record
go run ./cmd/bw overdub
go run ./cmd/bw panic
go run ./cmd/bw tempo 128
go run ./cmd/bw track create instrument 8
go run ./cmd/bw track 2 arm
go run ./cmd/bw scene create
go run ./cmd/bw scene 1 select
go run ./cmd/bw scene 1 launch
go run ./cmd/bw clip 4 4 create 8
go run ./cmd/bw clip 4 4 note 60 0 4
go run ./cmd/bw clip 4 4 pattern test-c
go run ./cmd/bw clip 2 4 launch
```

Use `--addr` before the command to target a non-default bridge address:

```powershell
go run ./cmd/bw --addr 127.0.0.1:58770 play
```

Restart Bitwig after rebuilding the Java extension:

```powershell
.\Restart-Bitwig.bat
```

The PowerShell script first asks Bitwig to close normally, waits 10 seconds, then
saves through the CLI bridge, force-stops remaining Bitwig processes if needed,
and launches Bitwig again. Use `.\Restart-Bitwig.bat -NoForce` to avoid the
force-stop fallback, or `.\Restart-Bitwig.bat -SkipSave` to skip the save step.

To restart directly into a project:

```powershell
.\Restart-Bitwig.bat -ProjectPath "C:\Users\jerem\OneDrive\Documents\Bitwig Studio\Projects\codex-performance-lab\codex-performance-lab\codex-performance-lab.bwproject"
```

For the Codex performance lab project:

```powershell
.\Start-Performance-Lab.bat
```

## Current Status

This is a first working scaffold, not a finished installed Bitwig extension.

- The Go CLI parses commands and sends JSON over UDP.
- The extension folder contains a Maven project for the Bitwig-side listener.
- The Bitwig-side source has been compiled, installed into Bitwig, and
  smoke-tested for basic transport commands.

## Design Notes

- UDP keeps the performance bridge low ceremony and local-only.
- Commands are named around live performance actions first.
- The CLI avoids writing project files or touching Bitwig project internals.
- `panic` is mapped to transport stop in the first version.
- Clip note steps use MIDI pitch numbers and 16th-note step positions.

## Next Steps

1. Build the extension from `bitwig-extension` with `mvn install`.
2. Install `target/BitwigCliBridge.bwextension` into Bitwig's controller
   extension location.
3. Start Bitwig, enable `Bitwig CLI Bridge`, and run `bw status`.
4. Expand commands around live looping: named tracks, next-empty-slot recording,
   scene snapshots, and macro morphs.
