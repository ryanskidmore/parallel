package parallel

import (
	"errors"
	"sync"
)

// The WorkerHelper struct contains components to
// assist the execution of workers.
type WorkerHelper struct {
	worker *Worker
	wg     *sync.WaitGroup
}

func newWorkerHelper(w *Worker) *WorkerHelper {
	wg := &sync.WaitGroup{}
	return &WorkerHelper{worker: w, wg: wg}
}

// Signal to the worker helper that this worker is
// complete. This must be run in each thread if you
// are calling Wait().
func (wh *WorkerHelper) Done() {
	wh.wg.Done()
}

// Publish data through the specified DataChannel.
// This is a non-blocking operation.
func (wh *WorkerHelper) PublishData(name string, data interface{}) error {
	if _, exists := wh.worker.p.dataChannels[name]; !exists {
		return errors.New("Data channel does not exist")
	}
	go func() {
		wh.worker.p.dataChannels[name] <- data
	}()
	return nil
}

// Consume data from the specified DataChannel.
// This is a blocking operation.
func (wh *WorkerHelper) ConsumeData(name string) (interface{}, error) {
	if _, exists := wh.worker.p.dataChannels[name]; !exists {
		return nil, errors.New("Data channel does not exist")
	}
	data, open := <-wh.worker.p.dataChannels[name]
	if !open {
		return nil, errors.New("Data channel closed")
	}
	return data, nil
}

// Consume data from the specified DataChannel
// in batches. This is a blocking operation and
// will only run when there are enough items to
// batch.
func (wh *WorkerHelper) ConsumeDataInBatches(name string, size int) ([]interface{}, error) {
	if _, exists := wh.worker.p.dataChannels[name]; !exists {
		return nil, errors.New("Data channel does not exist")
	}
	dataBatch := make([]interface{}, size, size)
	for i := 0; i < size; i++ {
		data, open := <-wh.worker.p.dataChannels[name]
		if !open {
			return dataBatch, errors.New("Data channel closed")
		}
		dataBatch[i] = data
	}
	return dataBatch, nil
}
