# termiflow-cli

Distilled tech news in your terminal. Powered by the [Termiflow](https://termiflow.com) managed backend — no API keys required.

> Bloomberg terminal aesthetic. Curated by AI. Built for developers.

## Install

```bash
go install github.com/oluoyefeso/termiflow-cli/cmd/termiflow@latest
```

Or build from source:

```bash
git clone https://github.com/oluoyefeso/termiflow-cli
cd termiflow-cli
go build -o termiflow ./cmd/termiflow
```

## Quick Start

```bash
# Configure with your Termiflow API key
termiflow config

# Subscribe to a topic
termiflow subscribe "rust async"
termiflow subscribe "LLM infrastructure"

# Read your feed
termiflow feed

# Continuous watch mode
termiflow feed --watch --interval 30m

# Ask a question with web search
termiflow ask "What changed in Rust 1.82?"
```

## Commands

| Command | Description |
|---------|-------------|
| `termiflow config` | Configure your API key |
| `termiflow subscribe <topic>` | Subscribe to a topic |
| `termiflow topics` | List your subscriptions |
| `termiflow unsubscribe <topic>` | Unsubscribe from a topic |
| `termiflow feed` | Read your curated feed |
| `termiflow feed --watch` | Continuous feed mode |
| `termiflow ask <question>` | Ask a question |
| `termiflow version` | Show version |

## Get an API Key

Sign up at [termiflow.com](https://termiflow.com) to get your `tf_xxx` API key.

## Contributing

### Mock mode

For development without real API calls, use mock mode:

```bash
TERMIFLOW_MOCK=true go run -tags mock ./cmd/termiflow feed --refresh
TERMIFLOW_MOCK=true go run -tags mock ./cmd/termiflow ask "test question"
```

Mock providers return canned responses so you can develop and test UI changes without
hitting the backend or spending API credits.

### Running tests

```bash
go test ./...                    # all tests
go test -race ./...              # with race detector
go test -v ./internal/cli/...    # specific package
```

## Self-hosting

Want to bring your own Anthropic and Tavily keys? Use the [self-hosted version](https://github.com/oluoyefeso/termiflow) instead.
