package streamserver

import "log"

type ConnLimiter struct {
	concurrentConn int
	bucket         chan int
}

func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.concurrentConn {
		log.Println("reached the rate limitation.")
		return false
	}

	cl.bucket <- 1
	return true
}

func (cl *ConnLimiter) ReleaseConn() {
	<-cl.bucket
}
