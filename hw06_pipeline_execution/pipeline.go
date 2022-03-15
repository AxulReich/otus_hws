package hw06pipelineexecution

import "sync"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var (
		out   Out
		wg    sync.WaitGroup
		ready = make(chan struct{})
	)

	go func() {
		for i := range stages {
			wg.Add(1)
			out = stages[i](controller(&wg, in, done))
			in = out
		}

		ready <- struct{}{}
		wg.Wait()
	}()

	<-ready

	return out
}
func controller(wg *sync.WaitGroup, in In, done In) Bi {
	out := make(Bi)

	go func() {
		defer wg.Done()
		defer close(out)

		for data := range in {
			select {
			case <-done:
				return
			default:
			}

			select {
			case <-done:
				return
			default:
				out <- data
			}
		}
	}()

	return out
}
