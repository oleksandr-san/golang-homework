package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func doneStage(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for range in {
			}
		}()

		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}

				select {
				case out <- v:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	prevStageOut := in
	for _, stage := range stages {
		prevStageOut = stage(doneStage(prevStageOut, done))
	}
	return doneStage(prevStageOut, done)
}
