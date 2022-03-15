package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var out Out

	for i := range stages {
		out = stages[i](middleWare(in, done))
		in = out
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
			case val, ok := <-in:
				if ok {
					out <- val
				} else {
					return
				}
			}
		}
	}()

	return out
}
