package parallel

import "errors"

// The Worker struct contains the necessary components for
// a parallelizable worker.
type Worker struct {
	p       *Parallel
	Name    string
	Config  *WorkerConfig
	execute func(wh *WorkerHelper, args interface{})
	helper  *WorkerHelper
}

// The WorkerConfig struct contains the configuration for
// a worker.
type WorkerConfig struct {
	Parallelism int
}

// Creates a new worker with the specified name and config.
func (p *Parallel) NewWorker(name string, cfg *WorkerConfig) (*Worker, error) {
	if _, exists := p.workers[name]; exists {
		return nil, errors.New("worker already exists")
	}
	if cfg.Parallelism < 1 {
		return nil, errors.New("parallelism must be 1 or higher")
	}
	w := &Worker{
		p:      p,
		Name:   name,
		Config: cfg,
	}
	p.workers[name] = w
	return w, nil
}

// Get a worker by name.
func (p *Parallel) Worker(name string) *Worker {
	if _, exists := p.workers[name]; !exists {
		return nil
	}
	return p.workers[name]
}

// Set the execution function of the worker. This is
// the function that is executed inside every worker.
func (w *Worker) SetExecution(exec func(wh *WorkerHelper, args interface{})) {
	w.execute = exec
}

// Start a worker with the specified args, which are
// passed to every instance of the worker.
func (w *Worker) Start(args interface{}) {
	wh := newWorkerHelper(w)
	w.helper = wh
	for i := 0; i < w.Config.Parallelism; i++ {
		w.helper.wg.Add(1)
		go w.execute(wh, args)
	}
}

// Wait until all worker routines have finished
// processing.
func (w *Worker) Wait() {
	w.helper.wg.Wait()
}

// Set the number of parallel routines for a worker.
func (w *Worker) SetParallelism(p int) {
	w.Config.Parallelism = p
}
