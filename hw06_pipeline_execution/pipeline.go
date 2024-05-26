package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outChannel := make(Bi)
	for _, stage := range stages {
		in = stage(in)
	}
	go func() {
		defer close(outChannel)
		for i := range in {
			select {
			case <-done:
				return
			case outChannel <- i:
			}
		}
	}()
	return outChannel
}
