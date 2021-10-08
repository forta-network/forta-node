package workers

// Workers receive work, execute and send results back.
type Workers struct {
	work chan *Work
}

// New creates new workers.
func New(size int) *Workers {
	workers := &Workers{
		work: make(chan *Work),
	}

	for i := 0; i < size; i++ {
		go workers.listenAndExecute()
	}

	return workers
}

func (workers *Workers) listenAndExecute() {
	for work := range workers.work {
		values, error := work.Func()
		work.OutputsCh <- &Output{
			Values: values,
			Error:  error,
		}
		close(work.OutputsCh)
	}
}

// Execute sends work to listening goroutines.
func (workers *Workers) Execute(f WorkFunc) *Output {
	outputsCh := make(chan *Output)
	workers.work <- &Work{
		Func:      f,
		OutputsCh: outputsCh,
	}
	return <-outputsCh
}

// WorkFunc is the work logic.
type WorkFunc func() ([]interface{}, error)

// Work defines the work to deal with.
type Work struct {
	Func      WorkFunc
	OutputsCh chan *Output
}

// Output defines the type for output values.
type Output struct {
	Values []interface{}
	Error  error
}
