package parallel

import "errors"

// Create a new DataChannel with specified name. A DataChannel
// is a conduit for transferring data between workers/goroutines.
func (p *Parallel) NewDataChannel(name string) error {
	if _, exists := p.dataChannels[name]; exists {
		return errors.New("Data channel already exists")
	}
	p.dataChannels[name] = make(chan interface{})
	return nil
}

// Close an existing DataChannel with specified name.
func (p *Parallel) CloseDataChannel(name string) error {
	if _, exists := p.dataChannels[name]; !exists {
		return errors.New("Data channel doesn't exists")
	}
	close(p.dataChannels[name])
	delete(p.dataChannels, name)
	return nil
}
