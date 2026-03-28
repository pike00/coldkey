# coldkey

Post-quantum age key generation and paper backup tool.

Generates [ML-KEM-768 + X25519](https://words.filippo.io/post-quantum-age/) hybrid age keys and produces single-page printable HTML backups with QR codes for disaster recovery.

## Quick start

### Docker (recommended)

```bash
# Pull the image
docker pull ghcr.io/pike00/coldkey:latest

# Interactive — generate a key and paper backup
just docker-run

# Backup an existing key
just docker-backup ~/.config/sops/age/keys.txt
```

All `just docker-*` commands include security hardening flags (network isolation, read-only filesystem, dropped capabilities). Output is written to `./output/`.

### From source

```bash
go install github.com/pike00/coldkey/cmd/coldkey@latest
coldkey generate -o ~/.config/sops/age/keys.txt
```

## Commands

### `coldkey` (no args) — Interactive mode

Presents a menu to generate a new key or create a backup from an existing one. Prompts for file paths and confirms before overwriting.

### `coldkey generate`

Generate a new post-quantum age key pair.

```
coldkey generate [flags]
  -o PATH       Key file output path (default: stdout)
  -f            Overwrite existing file
  --no-backup   Skip HTML backup generation
```

### `coldkey backup`

Create a printable HTML paper backup from an existing key file.

```
coldkey backup [flags] KEYFILE
  -o PATH    HTML output path (default: KEYFILE-backup.html)
```

### `coldkey version`

Print the version string.

## Security model

| Layer | Measure |
|-------|---------|
| Memory | `mlockall(MCL_CURRENT\|MCL_FUTURE)` prevents key material from being swapped to disk |
| Files | Written with mode `0600`, fsynced; temporaries shredded (3-pass overwrite) |
| Process | Secrets passed via stdin/files only, never in process arguments |
| Container | `--network none --read-only --cap-drop ALL --security-opt no-new-privileges:true` |
| Image | `distroless/static:nonroot` — no shell, non-root UID 65534 |
| Memory zeroing | Best-effort `secure.Zero()` on key buffers before GC (see [Limitations](#limitations)) |

### Docker flags explained

The `just docker-run` and `just docker-backup` commands apply these flags automatically:

| Flag | Purpose |
|------|---------|
| `--network none` | No network access — key generation is purely local |
| `--read-only` | Immutable root filesystem |
| `--cap-drop ALL` | Drop all Linux capabilities |
| `--security-opt no-new-privileges:true` | Prevent privilege escalation |
| `--tmpfs /tmp:rw,noexec,nosuid,size=10m` | RAM-backed temp directory |
| `--cap-add IPC_LOCK` | (Optional) Enable `mlockall` for swap protection |

## QR code encoding

PQ age stores only the 32-byte seed (not the expanded ML-KEM-768 private key), so the full `keys.txt` is typically ~2,089 bytes — fitting in a single QR code (version 40, EC-L supports 2,953 bytes).

If a key file exceeds single-QR capacity, coldkey automatically splits it across multiple QR codes using a simple framing protocol:

```
COLDKEY:<part>/<total>:<data>
```

Recovery: scan all QR codes in order, strip the `COLDKEY:N/M:` prefix from each, concatenate, and verify the SHA-256 checksum.

## Paper backup contents

The generated HTML document contains:

- Title and metadata (date, hostname, user, source path)
- Raw key text in monospace (for manual transcription)
- QR code(s) with capacity annotation
- SHA-256 checksum for verification
- Step-by-step recovery instructions
- Print button (hidden in print media)

## Recovery procedure

1. Scan the QR code (or type the raw key text)
2. Save to `~/.config/sops/age/keys.txt`
3. Verify: `sha256sum keys.txt` matches the printed checksum
4. Test: `sops -d <any .sops file>`

## Building

```bash
just build       # Local binary
just docker      # Docker image (ghcr.io/pike00/coldkey)
just test        # Run tests
just ci          # Full CI: vet → test → build → docker
```

## Limitations

- **Go GC and secure memory**: Go's garbage collector may copy objects in memory, and Go strings are immutable, meaning key material held as a `string` (e.g. from `identity.String()`) cannot be reliably overwritten. `secure.Zero()` uses Go's built-in `clear()` to erase `[]byte` buffers, but earlier string copies may persist in the heap until garbage collected. `mlockall` prevents any of this from being swapped to disk; together these provide defense-in-depth, not a cryptographic guarantee that key material is erased from RAM immediately.
- **`mlockall` requires `CAP_IPC_LOCK`**: Add `--cap-add IPC_LOCK` to Docker run for full swap protection. Without it, coldkey prints a warning to stderr and continues.
- **QR scanning**: Very dense QR codes (version 40) may be hard to scan from paper. The raw key text is always included as a manual fallback.

## License

MIT
