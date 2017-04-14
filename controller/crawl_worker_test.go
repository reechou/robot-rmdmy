package controller

import (
	"testing"
	"time"
)

func TestCrawlWorker(t *testing.T) {
	cw := NewCrawlWorker()
	time.Sleep(60 * time.Second)
	cw.Stop()
}
