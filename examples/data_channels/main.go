package main

import (
	"log"
	"time"

	"github.com/ryanskidmore/parallel"
)

func main() {
	p := parallel.New()
	err := p.NewDataChannel("time")
	if err != nil {
		log.Fatalf("FATAL: Failed to create new datachannel: %v", err)
	}
	wt, err := p.NewWorker("time", &parallel.WorkerConfig{Parallelism: 1})
	if err != nil {
		log.Fatalf("FATAL: Failed to create new worker: %v", err)
	}
	wt.SetExecution(func(wh *parallel.WorkerHelper, args interface{}) {
		for i := 0; i < 10; i++ {
			err := wh.PublishData("time", time.Now().Unix())
			if err != nil {
				log.Printf("ERROR: Failed to publish to datachannel: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
		wh.Done()
	})
	wp, err := p.NewWorker("time_printer", &parallel.WorkerConfig{Parallelism: 1})
	if err != nil {
		log.Fatalf("FATAL: Failed to create new worker: %v", err)
	}
	wp.SetExecution(func(wh *parallel.WorkerHelper, args interface{}) {
		for {
			data, err := wh.ConsumeData("time")
			if err != nil {
				break
			}
			log.Printf("INFO: Current time is %v", time.Unix(data.(int64), 0))
		}
		wh.Done()
	})
	wt.Start(nil)
	wp.Start(nil)
	wt.Wait()
	err = p.CloseDataChannel("time")
	if err != nil {
		log.Fatalf("FATAL: Failed to close datachannel: %v", err)
	}
	wp.Wait()
}
