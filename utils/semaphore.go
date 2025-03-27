package utils

type Semaphore struct {
	size     int
	resource chan struct{}
}

func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{
		size:     limit,
		resource: make(chan struct{}, limit)}
}

func (sem *Semaphore) Acquire() {
	sem.resource <- struct{}{}
}

func (sem *Semaphore) Release() {
	<-sem.resource
}

func (sem *Semaphore) Size() int {
	return sem.size
}

func (sem *Semaphore) Close() {
	close(sem.resource)
}
