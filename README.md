# Eino

[![Go Report Card](https://goreportcard.com/badge/github.com/cloudwego/eino)](https://goreportcard.com/report/github.com/cloudwego/eino)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](go.mod)

Eino is a Go framework for building LLM-powered applications with composable components and type-safe pipelines.

This repository is a fork of [cloudwego/eino](https://github.com/cloudwego/eino).

## Features

- **Composable Components**: Build complex AI workflows from reusable building blocks
- **Type-Safe Pipelines**: Leverage Go generics for compile-time safety
- **Streaming Support**: First-class support for streaming LLM responses
- **Graph-Based Orchestration**: Define complex multi-step workflows as directed graphs
- **Extensible**: Easy to add custom components, models, and tools

## Getting Started

### Prerequisites

- Go 1.21 or later

### Installation

```bash
go get github.com/cloudwego/eino
```

### Quick Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/cloudwego/eino/components/model"
    "github.com/cloudwego/eino/schema"
)

func main() {
    ctx := context.Background()

    // Initialize your LLM model
    // llm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
    //     Model:  "gpt-4o",
    //     APIKey: os.Getenv("OPENAI_API_KEY"),
    // })

    messages := []*schema.Message{
        schema.SystemMessage("You are a helpful assistant."),
        schema.UserMessage("Hello, how are you?"),
    }

    _ = messages
    _ = ctx
    fmt.Println("Eino is ready to use!")
    log.Println("See examples/ for complete usage")
}
```

## Project Structure

```
eino/
├── components/          # Core component interfaces
│   ├── document/        # Document loader and transformer interfaces
│   ├── embedding/       # Embedding model interfaces
│   ├── indexer/         # Vector store indexer interfaces
│   ├── model/           # LLM chat model interfaces
│   ├── prompt/          # Prompt template interfaces
│   ├── retriever/       # Retriever interfaces
│   └── tool/            # Tool/function calling interfaces
├── compose/             # Pipeline and graph composition
├── flow/                # Pre-built workflow patterns
├── schema/              # Core data types and schemas
└── utils/               # Shared utilities
```

## Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) and check the [issue tracker](https://github.com/cloudwego/eino/issues) for open tasks.

### Development

```bash
# Clone the repository
git clone https://github.com/cloudwego/eino.git
cd eino

# Run tests
go test ./...

# Run linter
golangci-lint run
```

## Personal Notes

> **Note (personal fork):** I'm using this fork primarily to experiment with custom retriever implementations and graph-based RAG pipelines. The `examples/` directory contains my own usage examples not present in the upstream repo.

## License

This project is licensed under the Apache 2.0 License — see the [LICENSE](LICENSE) file for details.

## Acknowledgements

This project is a fork of [cloudwego/eino](https://github.com/cloudwego/eino), originally developed by the CloudWeGo team.
