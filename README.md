# parallel: a Go Parallel Processing Library

Concurrency is hard. This library doesn't aim to make it easy, but it will hopefully make it a little less painful. 

## Install

`go get github.com/ryanskidmore/parallel`

## Testing

`go test`

## Usage

Initialise an instance of parallel

`p := parallel.New()`

Create a new worker. A worker is a processor that runs inside a goroutine. The parallelism configuration option defines how many goroutines running this worker will be spun up. 

```go
worker, err := p.NewWorker("name", &parallel.WorkerConfig{Parallelism: 1})
```

Retrieve a worker.

```go
worker := p.Worker("name")
```

Set the execution function. The execution function is the function that is executed inside the goroutine. This function is only run once per goroutine. 

The execution function must have the signature `func(wh *parallel.WorkerHelper, args interface{})`.

```go
worker.SetExecution(func(wh *parallel.WorkerHelper, args interface{}) { fmt.Println(args) })
```

Starting the worker. This will start the worker in the configured number of goroutines with the given args.

```go
worker.Start(interface{}("args"))
```

Wait for the worker to complete. This function waits until goroutines have finished executing. Note: you must call `WorkerHelper.Done()` in your execution function for this to work

```go
worker.Wait()
```

### WorkerHelper

The WorkerHelper struct gives you access to a WaitGroup via `WorkerHelper.Done()`.

The WorkerHelper can also be used to consume and publish data from other workers/goroutines. This is done using DataChannels.

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