package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outChannel := in
	taskBinding := func(in In, done In) Out {
		out := make(Bi)

		go func() {
			defer close(out)
			for {
				select {
				case <-done:
					return
				case i, ok := <-in:
					if !ok {
						return
					}
					out <- i
				}
			}
		}()
		return out
	}

	var count int
	for i := range stages {
		count++
		outChannel = taskBinding(stages[i](outChannel), done)
	}
	return outChannel
}
