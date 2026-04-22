# ccmgr

Claude Code Configuration Manager — a TUI tool to manage multiple environment configurations for [Claude Code](https://docs.anthropic.com/en/docs/claude-code).

Switch between providers like Z.AI and default Anthropic with a single command.

## Installation

### Download binary (recommended)

Grab the latest release for your platform from [Releases](https://github.com/smaddiona/ccmgr/releases).

```bash
# macOS (Apple Silicon)
tar xzf ccmgr_*_darwin_aarch64.tar.gz
chmod +x ccmgr
sudo mv ccmgr /usr/local/bin/

# macOS (Intel)
tar xzf ccmgr_*_darwin_x86_64.tar.gz
chmod +x ccmgr
sudo mv ccmgr /usr/local/bin/

# Linux (amd64)
tar xzf ccmgr_*_linux_x86_64.tar.gz
chmod +x ccmgr
sudo mv ccmgr /usr/local/bin/
```

### Go install

```bash
go install github.com/smaddiona/ccmgr@latest
```

### Build from source

```bash
git clone https://github.com/smaddiona/ccmgr.git
cd ccmgr
go build -o ccmgr .
```

## Usage

### Create a profile

```bash
ccmgr create
```

Interactive TUI wizard that walks you through:

1. **Select preset** — Choose a provider (Z.AI, Default Anthropic, etc.)
2. **Enter API key** — Paste your provider API key (masked input)
3. **Assign models** — Set Opus, Sonnet, and Haiku model names (pre-filled with defaults for known providers)
4. **Name the profile** — Give it a label like "Z.AI Production"

For the **Default** preset, steps 2 and 3 are skipped since no env vars are needed.

### Switch profile

```bash
ccmgr switch
```

Shows a TUI list of all profiles. Select one to activate — the tool writes the env vars to `~/.claude/settings.json` and preserves all other settings.

### List profiles

```bash
ccmgr list
```

Plain text output showing all profiles. The active profile is marked with `*`.

### Delete a profile

```bash
ccmgr delete          # TUI selector
ccmgr delete "My Profile"  # direct by label
```

Removes the profile from storage. If the deleted profile was active, the env section in `~/.claude/settings.json` is cleared.

## How it works

```
~/.ccmgr/
  └── profiles.json    # stored profiles (0600 permissions)

~/.claude/
  └── settings.json    # Claude Code config (env section managed by ccmgr)
```

- **Profiles** are stored separately in `~/.ccmgr/profiles.json`
- Only the `env` section of `~/.claude/settings.json` is modified — all other settings (`model`, `enabledPlugins`, etc.) are preserved
- A backup is created at `~/.claude/settings.json.bak` before each write

### Example profile

A Z.AI profile produces this env section in `settings.json`:

```json
{
  "env": {
    "ANTHROPIC_AUTH_TOKEN": "<your-api-key>",
    "ANTHROPIC_BASE_URL": "https://api.z.ai/api/anthropic",
    "ANTHROPIC_DEFAULT_HAIKU_MODEL": "glm-4.5-air",
    "ANTHROPIC_DEFAULT_OPUS_MODEL": "glm-5.1",
    "ANTHROPIC_DEFAULT_SONNET_MODEL": "glm-4.7",
    "API_TIMEOUT_MS": "3000000"
  }
}
```

Switching to the **Default** preset removes the `env` section entirely.

## Supported presets

| Preset | Base URL | API Key | Models |
|--------|----------|---------|--------|
| Default | — | No | No |
| Z.AI | `https://api.z.ai/api/anthropic` | Yes | glm-5.1, glm-4.7, glm-4.5-air |

## Development

```bash
go test ./...
go build -o ccmgr .
```

## License

MIT
