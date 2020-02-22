package parallel

import (
	"sync"
	"testing"
)

func TestWorker(t *testing.T) {
	t.Run("Initialization", func(t *testing.T) {
		p := New()
		t.Run("NewWorkerIncorrectConfig", func(t *testing.T) {
			_, err := p.NewWorker("TestWorker", &WorkerConfig{Parallelism: 0})
			testNotNil(t, err)
		})
		t.Run("NewWorker", func(t *testing.T) {
			w, err := p.NewWorker("TestWorker", &WorkerConfig{Parallelism: 1})
			testNil(t, err)
			testNotNil(t, w)
		})
		t.Run("DuplicateNewWorker", func(t *testing.T) {
			_, err := p.NewWorker("TestWorker", &WorkerConfig{Parallelism: 1})
			testNotNil(t, err)
		})
		t.Run("GetWorker", func(t *testing.T) {
			w := p.Worker("TestWorker")
			testAssert(t, w.Name == "TestWorker", "GetWorker failed: incorrect name")
			testNotNil(t, w.Config)
			testAssert(t, w.Config.Parallelism == 1, "GetWorker failed: incorrect parallelism")
		})
		t.Run("GetWorkerDoesntExist", func(t *testing.T) {
			w := p.Worker("DoesntExist")
			testAssert(t, w == nil, "worker is not nil")
		})
		t.Run("SetParallelism", func(t *testing.T) {
			p.Worker("TestWorker").SetParallelism(12)
			testAssert(t, p.Worker("TestWorker").Config.Parallelism == 12, "GetWorker failed: incorrect parallelism")
		})
	})
	t.Run("Execution", func(t *testing.T) {
		p := New()
		w, err := p.NewWorker("TestWorker", &WorkerConfig{Parallelism: 1})
		if err != nil {
			t.Fatalf("ExecutionSetup: Failed to create new worker: %v", err)
		}
		t.Run("ExecutionSingle", func(t *testing.T) {
			ts := &testStruct{Counter: 0, Mutex: sync.Mutex{}}
			ef := func(wh *WorkerHelper, args interface{}) {
				ts := args.(*testStruct)
				ts.Mutex.Lock()
				ts.Counter = ts.Counter + 1
				ts.Mutex.Unlock()
				wh.Done()
			}
			w.SetExecution(ef)
			w.Start(interface{}(ts))
			w.Wait()
			testAssert(t, ts.Counter == 1, "Execution failed: counter does not equal 1 (equals %v)", ts.Counter)
		})
		t.Run("ExecutionMultiple", func(t *testing.T) {
			ts := &testStruct{Counter: 0, Mutex: sync.Mutex{}}
			ef := func(wh *WorkerHelper, args interface{}) {
				ts := args.(*testStruct)
				ts.Mutex.Lock()
				ts.Counter = ts.Counter + 1
				ts.Mutex.Unlock()
				wh.Done()
			}
			w.SetExecution(ef)
			w.SetParallelism(8)
			w.Start(interface{}(ts))
			w.Wait()
			testAssert(t, ts.Counter == 8, "Execution failed: counter does not equal 8 (equals %v)", ts.Counter)
		})
	})
}
