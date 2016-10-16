package devproxy

import (
	"strconv"
	"sync"
	"testing"
)

func TestIdGeneratorHandler(t *testing.T) {
	idGenerator := NewStringIdGenerator()
	for i := 1; i <= 10; i++ {
		if idGenerator.NewId() != strconv.Itoa(i) {
			t.Errorf("Expected %d", i)
		}
	}
}

func TestIdGeneratorHandlerWithMultipleThreads(t *testing.T) {
	idGenerator := NewStringIdGenerator()
	result := make(map[string]bool)
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			// TODO: use channel!?
			result[idGenerator.NewId()] = true
		}()
	}
	wg.Wait()
	for i := 1; i <= 100; i++ {
		if !result[strconv.Itoa(i)] {
			t.Errorf("%d not set", i)
		}
	}
}
