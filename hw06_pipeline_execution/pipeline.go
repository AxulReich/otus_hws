package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {

	go func() {
		var (
			wg sync.WaitGroup
		)
		for _, stage := range stages {
			wg.Add(1)
			stage := stage
			go func() {
				out := stage(in)
				in = outStage
			}()

		}

		wg.Wait()
	}()

	return out
}
