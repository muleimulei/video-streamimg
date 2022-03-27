package taskrunner

import "time"

type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

func NewWorker(interval time.Duration, r *Runner) *Worker {
	return &Worker{
		ticker: time.NewTicker(interval),
		runner: r,
	}
}

func (w *Worker) startWorker() {
	for {
		select {
		case <-w.ticker.C:
			go w.runner.StartAll()
		}
	}
}

func Start() {
	// start video file cleaning

	r := NewRunner(3, true, VideoClearDispatcher, VideoClearExecutor)
	w := NewWorker(time.Minute, r)

	go w.startWorker()
}
