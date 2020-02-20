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
		Test_NotNil(t, wh)
		Test_NotNil(t, wh.worker)
		Test_NotNil(t, wh.wg)
	})
	t.Run("PublishData", func(t *testing.T) {
		wh := newWorkerHelper(w)
		err := wh.PublishData("TestChannel", "testdata")
		Test_Nil(t, err)
	})
	t.Run("ConsumeData", func(t *testing.T) {
		wh := newWorkerHelper(w)
		data, err := wh.ConsumeData("TestChannel")
		Test_Nil(t, err)
		Test_Equals(t, data.(string), "testdata")
	})
	t.Run("ConsumeDataInBatches", func(t *testing.T) {
		wh := newWorkerHelper(w)
		for i := 0; i < 100; i++ {
			err := wh.PublishData("TestChannel", "testdata")
			Test_Nil(t, err)
		}
		for i := 0; i < 5; i++ {
			data, err := wh.ConsumeDataInBatches("TestChannel", 20)
			Test_Nil(t, err)
			Test_Assert(t, len(data) == 20, "length of batch isn't 20 (value: %v)", len(data))
		}
	})
	t.Run("ConsumeDataCloseChan", func(t *testing.T) {
		wh := newWorkerHelper(w)
		go func() {
			time.Sleep(50 * time.Millisecond)
			err = p.CloseDataChannel("TestChannel")
			Test_Nil(t, err)
		}()
		_, err := wh.ConsumeData("TestChannel")
		Test_NotNil(t, err)
	})
}
