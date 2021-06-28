package agentpool

import (
	"sync"
	"sync/atomic"
)

var processingState pauser

type pauser struct {
	paused int64
	wg     sync.WaitGroup
}

// Pause increments the wait group counter only once.
func (p *pauser) Pause() {
	canPause := atomic.CompareAndSwapInt64(&p.paused, 0, 1)
	if canPause {
		p.wg.Add(1)
	}
}

// Continue decrements the wait group counter only once.
func (p *pauser) Continue() {
	canContinue := atomic.CompareAndSwapInt64(&p.paused, 1, 0)
	if canContinue {
		p.wg.Done()
	}
}

func (p *pauser) waitIfPaused() {
	p.wg.Wait()
}
