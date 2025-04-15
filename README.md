# Signals

The `signals` a robust, dependency-free go library that provides simple, thin, and user-friendly pub-sub kind of in-process event system for your Go applications. It allows you to generate and emit signals (synchronously or asynchronously) as well as manage listeners.

💯 **100% test coverage** 💯

[![GoReportCard example](https://goreportcard.com/badge/github.com/nanomsg/mangos)](https://goreportcard.com/report/github.com/maniartech/signals)
[![<ManiarTech®️>](https://circleci.com/gh/maniartech/signals.svg?style=shield)](https://circleci.com/gh/maniartech/signals)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/maniartech/signals)

## Installation

```bash
go get github.com/maniartech/signals
```

## Usage

```go
package main

import (
  "context"
  "fmt"
  
  "github.com/linux019/signals"
)

var RecordCreated = signals.New[Record]()
var RecordUpdated = signals.New[Record]()
var RecordDeleted = signals.New[Record]()

func main() {

  // Add a listener to the RecordCreated signal
  RecordCreated.AddListener(func(ctx context.Context, record Record) {
    fmt.Println("Record created:", record)
  }, signals.SignalType(123)) // <- Key is optional useful for removing the listener later

  // Add a listener to the RecordUpdated signal
  RecordUpdated.AddListener(func(ctx context.Context, record Record) {
    fmt.Println("Record updated:", record)
  })

  // Add a listener to the RecordDeleted signal
  RecordDeleted.AddListener(func(ctx context.Context, record Record) {
    fmt.Println("Record deleted:", record)
  })

  ctx := context.Background()

  // Emit the RecordCreated signal
  RecordCreated.Emit(ctx, Record{ID: 1, Name: "John"})

  // Emit the RecordUpdated signal
  RecordUpdated.Emit(ctx, Record{ID: 1, Name: "John Doe"})

  // Emit the RecordDeleted signal
  RecordDeleted.Emit(ctx, Record{ID: 1, Name: "John Doe"})
}
```

## Documentation

[![GoDoc](https://godoc.org/github.com/maniartech/signals?status.svg)](https://godoc.org/github.com/maniartech/signals)

## License

![License](https://img.shields.io/badge/license-MIT-blue.svg)
