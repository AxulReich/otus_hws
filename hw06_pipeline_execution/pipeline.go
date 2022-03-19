package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for i := range stages {
		out = stages[i](middleWare(out, done))
	}

	return out
}

func middleWare(in In, done In) Bi {
	out := make(Bi)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			default:
			}

			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				out <- val
			}
		}
	}()

	return out
}
