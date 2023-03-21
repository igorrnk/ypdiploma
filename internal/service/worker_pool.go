package service

type Job struct {
	orderId int
}

func (job *Job) Do() {

}

type WorkerPool struct {
	jobCh chan *Job
}

func NewWorkerPool() *WorkerPool {
	pool := &WorkerPool{
		jobCh: make(chan *Job),
	}
	return pool
}

func (pool *WorkerPool) Run(count int) {
	for i := 0; i < count; i++ {
		go func() {
			for job := range pool.jobCh {
				job.Do()
			}
		}()
	}
}

func (pool *WorkerPool) Stop() {
	close(pool.jobCh)
}
