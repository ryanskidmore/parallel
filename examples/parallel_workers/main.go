package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/ryanskidmore/parallel"
)

func main() {
	p := parallel.New()
	err := p.NewDataChannel("numbers")
	if err != nil {
		log.Fatalf("FATAL: Failed to create new datachannel: %v", err)
	}
	r, err := p.NewWorker("random", &parallel.WorkerConfig{Parallelism: 4})
	if err != nil {
		log.Fatalf("FATAL: Failed to create new worker: %v", err)
	}
	r.SetExecution(func(wh *parallel.WorkerHelper, args interface{}) {
		for i := 0; i < 100; i++ {
			rInt := rand.Intn(100)
			err := wh.PublishData("numbers", rInt)
			if err != nil {
				log.Printf("ERROR: Failed to publish to datachannel: %v", err)
			}
		}
		wh.Done()
	})
	ra, err := p.NewWorker("random_average", &parallel.WorkerConfig{Parallelism: 1})
	if err != nil {
		log.Fatalf("FATAL: Failed to create new worker: %v", err)
	}
	ra.SetExecution(func(wh *parallel.WorkerHelper, args interface{}) {
		for {
			data, err := wh.ConsumeDataInBatches("numbers", 10)
			if err != nil {
				break
			}
			total := 0
			for _, rInt := range data {
				total = total + rInt.(int)
			}
			log.Printf("INFO: Average is %v", total/10)
		}
		wh.Done()
	})
	r.Start(nil)
	ra.Start(nil)
	r.Wait()
	time.Sleep(1 * time.Second)
	err = p.CloseDataChannel("numbers")
	if err != nil {
		log.Fatalf("FATAL: Failed to close datachannel: %v", err)
	}
	ra.Wait()
}
