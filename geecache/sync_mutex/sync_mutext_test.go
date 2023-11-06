package sync_mutex

import (
	"testing"
	"time"
)

func TestPrintOnce(t *testing.T) {
	for i := 0; i < 10; i++ {
		go printOnce(100)
	}
	time.Sleep(time.Second)
}

func TestPrintOnce2(t *testing.T) {
	for i := 0; i < 10; i++ {
		go printOnce2(100)
	}
	time.Sleep(time.Second)
}

func TestPrintOnce3(t *testing.T) {
	for i := 0; i < 10; i++ {
		go printOnce3(100)
	}
	time.Sleep(time.Second)
}
