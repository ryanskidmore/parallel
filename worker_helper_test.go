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
		test_NotNil(t, wh)
		test_NotNil(t, wh.worker)
		test_NotNil(t, wh.wg)
	})
	t.Run("PublishData", func(t *testing.T) {
		wh := newWorkerHelper(w)
		err := wh.PublishData("TestChannel", "testdata")
		test_Nil(t, err)
	})
	t.Run("ConsumeData", func(t *testing.T) {
		wh := newWorkerHelper(w)
		data, err := wh.ConsumeData("TestChannel")
		test_Nil(t, err)
		test_Equals(t, data.(string), "testdata")
	})
	t.Run("ConsumeDataInBatches", func(t *testing.T) {
		wh := newWorkerHelper(w)
		for i := 0; i < 100; i++ {
			err := wh.PublishData("TestChannel", "testdata")
			test_Nil(t, err)
		}
		for i := 0; i < 5; i++ {
			data, err := wh.ConsumeDataInBatches("TestChannel", 20)
			test_Nil(t, err)
			test_Assert(t, len(data) == 20, "length of batch isn't 20 (value: %v)", len(data))
		}
	})
	t.Run("ConsumeDataCloseChan", func(t *testing.T) {
		wh := newWorkerHelper(w)
		go func() {
			time.Sleep(50 * time.Millisecond)
			err = p.CloseDataChannel("TestChannel")
			test_Nil(t, err)
		}()
		_, err := wh.ConsumeData("TestChannel")
		test_NotNil(t, err)
	})
}
