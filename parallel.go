package parallel

// Parallel struct contains Workers and DataChannels.
type Parallel struct {
	workers      map[string]*Worker
	dataChannels map[string]chan interface{}
}

// New creates a new instance of the parallel struct.
func New() *Parallel {
	workers := make(map[string]*Worker)
	dataChannels := make(map[string]chan interface{})
	return &Parallel{workers: workers, dataChannels: dataChannels}
}
