package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	taskBinding := func(in, done In) Out {
		out := make(Bi)
		go func() {
			for {
				select {
				case <-done:
					close(out)
					for range in {
					}
					return
				case i, ok := <-in:
					if !ok {
						close(out)
						return
					}
					out <- i
				}
			}
		}()
		return out
	}

	outChannel := taskBinding(in, done)
	for i := range stages {
		outChannel = taskBinding(stages[i](outChannel), done)
	}
	return outChannel
}
