package devproxy

import (
	"strconv"
	"sync/atomic"
)

type StringIdGenerator interface {
	NewId() string
}

type stringIdGeneratorImpl struct {
	counter uint64
}

func (s *stringIdGeneratorImpl) NewId() string {
	return strconv.FormatUint(atomic.AddUint64(&s.counter, 1), 10)
}

func NewStringIdGenerator() StringIdGenerator {
	return &stringIdGeneratorImpl{0}
}
