package parallel

import (
	"testing"
	"time"
)

func TestWorkerHelper(t *testing.T) {
	p := New()
	w, err := p.NewWorker("TestWorker", &WorkerConfig{Parallelism: 1})
	if err != nil {
		t.Fatalf("Failed to create new worker: %v", err)
	}
	err = p.NewDataChannel("TestChannel")
	if err != nil {
		t.Fatalf("Failed to create new data channel: %v", err)
	}
	t.Run("Initialization", func(t *testing.T) {
		wh := newWorkerHelper(w)
		testNotNil(t, wh)
		testNotNil(t, wh.worker)
		testNotNil(t, wh.wg)
	})
	t.Run("PublishData", func(t *testing.T) {
		wh := newWorkerHelper(w)
		err := wh.PublishData("TestChannel", "testdata")
		testNil(t, err)
	})
	t.Run("PublishDataChanDoesntExist", func(t *testing.T) {
		wh := newWorkerHelper(w)
		err := wh.PublishData("DoesntExist", "testdata")
		testNotNil(t, err)
	})
	t.Run("ConsumeData", func(t *testing.T) {
		wh := newWorkerHelper(w)
		data, err := wh.ConsumeData("TestChannel")
		testNil(t, err)
		testEquals(t, data.(string), "testdata")
	})
	t.Run("ConsumeDataDoesntExist", func(t *testing.T) {
		wh := newWorkerHelper(w)
		_, err := wh.ConsumeData("DoesntExist")
		testNotNil(t, err)
	})
	t.Run("ConsumeDataInBatches", func(t *testing.T) {
		wh := newWorkerHelper(w)
		for i := 0; i < 100; i++ {
			err := wh.PublishData("TestChannel", "testdata")
			testNil(t, err)
		}
		for i := 0; i < 5; i++ {
			data, err := wh.ConsumeDataInBatches("TestChannel", 20)
			testNil(t, err)
			testAssert(t, len(data) == 20, "length of batch isn't 20 (value: %v)", len(data))
		}
	})
	t.Run("ConsumeDataInBatchesDoesntExist", func(t *testing.T) {
		wh := newWorkerHelper(w)
		for i := 0; i < 5; i++ {
			_, err := wh.ConsumeDataInBatches("DoesntExist", 20)
			testNotNil(t, err)
		}
	})
	t.Run("ConsumeDataCloseChan", func(t *testing.T) {
		wh := newWorkerHelper(w)
		go func() {
			time.Sleep(50 * time.Millisecond)
			err = p.CloseDataChannel("TestChannel")
			testNil(t, err)
		}()
		_, err := wh.ConsumeData("TestChannel")
		testNotNil(t, err)
	})
	t.Run("ConsumeDataInBatchesCloseChan", func(t *testing.T) {
		wh := newWorkerHelper(w)
		go func() {
			time.Sleep(50 * time.Millisecond)
			err = p.CloseDataChannel("TestChannel")
			testNil(t, err)
		}()
		_, err := wh.ConsumeDataInBatches("TestChannel", 20)
		testNotNil(t, err)
	})
}
