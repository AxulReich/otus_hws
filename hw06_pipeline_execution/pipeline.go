package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var out Out

	var (
		wg    sync.WaitGroup
		ready = make(chan struct{})
	)

	go func() {
		defer func() {
			ready <- struct{}{}
		}()

		for _, stage := range stages {
			wg.Add(1)
			stage := stage
			go func() {
				defer wg.Done()
				for data := range in {

				}

				out = stage(in)
				in = out
			}()
		}

		wg.Wait()
	}()

	<-ready

	return out
}
