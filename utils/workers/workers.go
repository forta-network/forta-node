package workers

// Group receives work, execute and send results back.
type Group struct {
	in chan *Work
}

// New creates new worker group.
func New(size int) *Group {
	group := &Group{
		in: make(chan *Work),
	}

	for i := 0; i < size; i++ {
		go group.listenAndExecute()
	}

	return group
}

func (group *Group) listenAndExecute() {
	for work := range group.in {
		values, error := work.Func()
		work.OutputsCh <- &Output{
			Values: values,
			Error:  error,
		}
		close(work.OutputsCh)
	}
}

// Execute sends work to listening goroutines.
func (group *Group) Execute(f WorkFunc) *Output {
	outputsCh := make(chan *Output)
	group.in <- &Work{
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
