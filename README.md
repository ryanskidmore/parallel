# parallel: a Go Parallel Processing Library

[![Build Status](https://travis-ci.org/ryanskidmore/parallel.svg)](https://travis-ci.org/ryanskidmore/parallel)
[![codecov](https://codecov.io/gh/ryanskidmore/parallel/branch/master/graph/badge.svg)](https://codecov.io/gh/ryanskidmore/parallel)
[![Go Report Card](https://goreportcard.com/badge/github.com/ryanskidmore/parallel)](https://goreportcard.com/report/github.com/ryanskidmore/parallel)
[![GoDoc](https://godoc.org/github.com/ryanskidmore/parallel?status.svg)](https://godoc.org/github.com/ryanskidmore/parallel)

Concurrency is hard. This library doesn't aim to make it easy, but it will hopefully make it a little less painful. 

## Install
This library _should_ be compatible with all recent and future versions of Go, and has no third party dependencies. 
```sh
go get -u github.com/ryanskidmore/parallel
```

You can then import the library

```go
import "github.com/ryanskidmore/parallel"
```

## Testing
This library uses the standard Go testing tools, and doesn't use any third party testing libraries.
```sh
go test
```

## Quick Start

```go
package main

import (
    "log"
    "fmt"

    "github.com/ryanskidmore/parallel"
)

func main() {
    p := parallel.New() // Create a new instance of parallel
    worker, err := p.NewWorker("worker1", &parallel.WorkerConfig{Parallelism: 1}) // Create a new worker
    if err != nil {
        log.Fatalf("FATAL: Failed to create new worker: %v", err)
    }
    worker.SetExecution(func(wh *parallel.WorkerHelper, args interface{}) { 
        fmt.Println(args)
        wh.Done()
    })
    worker.Start(interface{}("Test String"))
    worker.Wait()
}
```
## Docs
The best source of reference is the [GoDocs](https://godoc.org/github.com/ryanskidmore/parallel) for this library. 
Noted below are parts of the library that may not be immediately obvious from the docs or otherwise.


### WorkerHelper

The WorkerHelper struct gives you access to a WaitGroup via `WorkerHelper.Done()`. When you call `Worker.Wait()` this waits
on the WorkGroup for that worker and will block until every instance of the worker has called `WorkerHelper.Done()`.

The WorkerHelper can also be used to consume and publish data to/from other workers/goroutines and this is done using DataChannels.

### DataChannels

DataChannels are intended to be a method of publishing and consuming data between different workers and goroutines.

Before using a DataChannel in an execution function, it must first be initialised using:
```go
err := p.NewDataChannel("name")
```

Once the DataChannel has been initialised, it can be published to by calling:
```go
err := WorkerHelper.PublishData("name", data)
```
This is an asynchronous operation, so will not block execution.

Data can be consumed from the DataChannel via either calling:
```go
data, err := WorkerHelper.ConsumeData("name")
```
or
```go
data, err := WorkerHelper.ConsumeDataInBatches("name", 20)
```

These functions will return an error when the DataChannel doesn't exist or when the channel is closed.

## Examples

Examples of usage can be found in the [examples](https://github.com/ryanskidmore/parallel/tree/master/examples) directory.