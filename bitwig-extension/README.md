# Bitwig CLI Bridge Extension

This folder contains the bridge-side scaffold for a Bitwig Controller Extension.
The Go CLI sends JSON commands over UDP to `127.0.0.1:58770`; the extension
receives those commands inside Bitwig and maps them to transport, track, and
clip launcher actions.

Bitwig's Controller API documentation ships inside Bitwig Studio under:

`Help > Documentation > Developer Resources`

This folder is a Maven controller-extension project based on Bitwig's installed
Java controller template.

Build:

```powershell
mvn install
```

Maven must be installed and available on `PATH`. The project resolves Bitwig's
controller extension API from Bitwig's Maven repository.

Output:

```text
target/BitwigCliBridge.bwextension
```

Initial command contract:

```json
{"action":"play"}
{"action":"stop"}
{"action":"record"}
{"action":"overdub"}
{"action":"panic"}
{"action":"tempo","value":128}
{"action":"track.create.instrument","track":8}
{"action":"scene.create"}
{"action":"scene.select","slot":1}
{"action":"scene.launch","slot":1}
{"action":"track.arm","track":2}
{"action":"clip.create","track":4,"slot":4,"value":8}
{"action":"clip.note","track":4,"slot":4,"pitch":60,"step":0,"length":4}
{"action":"clip.pattern","track":4,"slot":4,"name":"test-c"}
{"action":"clip.launch","track":2,"slot":4}
```
