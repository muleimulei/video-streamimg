package taskrunner

import (
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	d := func(dc dataChan) error {
		for i := 0; i < 30; i++ {
			dc <- i
			log.Println("Dispatch send ", i)
		}
		return nil
	}

	e := func(dc dataChan) error {
	forloop:
		for {
			select {
			case s := <-dc:
				log.Println("Execute received ", s)
			default:
				break forloop
			}
		}
		return nil
	}

	runner := NewRunner(30, false, d, e)
	go runner.StartAll()
	time.Sleep(3 * time.Second)
}
